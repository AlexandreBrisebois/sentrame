#!/usr/bin/env python3
"""Exercise a disposable image through its public HTTP/OTLP and local volume.

Requires a working Docker daemon. Creates and removes only uniquely named test
containers/bind folders. It never uses the repository's review database or credentials.
"""
import http.client
import http.cookiejar
import json
import os
import secrets
import socket
import subprocess
import tempfile
import time
import urllib.error
import urllib.request
from pathlib import Path


def docker(*args):
    return subprocess.check_output(["docker", *args], text=True).strip()


def wait_for_ready(request, timeout=90, pause=time.sleep, clock=time.monotonic):
    deadline = clock() + timeout
    while clock() < deadline:
        try:
            with request("/ready") as response:
                if response.status == 200:
                    return
        except urllib.error.HTTPError as error:
            # A wrong route/auth response is not a startup race.
            if error.code != 503:
                raise
            error.close()
        except (urllib.error.URLError, TimeoutError, ConnectionError, http.client.RemoteDisconnected):
            pass
        pause(1)
    raise RuntimeError("Container readiness did not pass within 90 seconds")


def smoke_payload(suffix, server_now):
    # Keep the synthetic occurrence inside the receiver's event-time window.
    # Host/VM clocks are independent; future source time is correctly excluded
    # from a historical query ending at the receiver's current time.
    return {"resourceLogs": [{"resource": {"attributes": [
            {"key": "service.name", "value": {"stringValue": "sentrame-container-smoke"}}
        ]}, "scopeLogs": [{"logRecords": [{"timeUnixNano": str(int(server_now) - 1_000_000_000),
            "eventName": "container.smoke." + suffix, "severityNumber": 9}]}]}]}


def main():
    image = os.environ.get("IMAGE_UNDER_TEST", "sentrame:review")
    platform = os.environ.get("PLATFORM_UNDER_TEST", "linux/arm64")
    if platform not in ("linux/amd64", "linux/arm64"):
        raise SystemExit("Choose linux/amd64 or linux/arm64")
    suffix = secrets.token_hex(6)
    name = "sentrame-smoke-" + suffix
    data_temporary = tempfile.TemporaryDirectory(prefix="sentrame-smoke-data-")
    data_folder = Path(data_temporary.name) / "data"
    data_folder.mkdir()
    # Only this disposable, synthetic test mount is world-writable. Its parent
    # stays mode0700. Deployment uses the chosen UID/GID10001-owned mode0700 folder.
    # Container UID10001 can write; the host test user can unlink files at cleanup.
    data_folder.chmod(0o777)
    with socket.socket() as listener:
        listener.bind(("127.0.0.1", 0))
        port = listener.getsockname()[1]
    origin = f"http://127.0.0.1:{port}"
    owner=secrets.token_hex(24)
    csrf=""
    app_id=""
    demo_id=""
    volume=name+"-data"
    docker("volume","create",volume)
    jar = http.cookiejar.CookieJar()
    client = urllib.request.build_opener(urllib.request.HTTPCookieProcessor(jar))

    def request(path, body=None, bearer=None):
        headers = {"Origin": origin,"X-CSRF-Token":csrf}
        if body is not None:
            headers["Content-Type"] = "application/json"
        if bearer:
            headers["Authorization"] = "Bearer " + bearer
        req = urllib.request.Request(origin + path, headers=headers,
                                     data=None if body is None else json.dumps(body).encode())
        return client.open(req, timeout=10)

    def login():
        nonlocal csrf
        with request("/api/login", {"passphrase": owner}) as response:
            assert response.status == 200
            csrf=json.load(response)["csrf"]

    def snapshot_has_record():
        with request("/api/apps/"+app_id+"/snapshot") as response:
            data = json.load(response)
        matching = [e for e in data["discovered"] if e["name"] == "container.smoke." + suffix]
        assert len(matching) == 1 and matching[0]["count"] == 1, "Committed evidence missing or duplicated"
        assert any(e["name"] == "container.smoke." + suffix for e in data["model"]["unassigned"]), "Standard record must remain accepted-unassigned"

    try:
        with tempfile.TemporaryDirectory(prefix="sentrame-container-env-") as temporary:
            env_file = Path(temporary) / "runtime.env"
            env_file.write_text("\n".join([
                "SENTRAME_ORIGIN=" + origin, "SENTRAME_PASSPHRASE=" + owner,
            ]) + "\n")
            env_file.chmod(0o600)
            docker("run", "--detach", "--name", name, "--platform", platform,
                   "--read-only", "--cap-drop", "ALL", "--security-opt", "no-new-privileges:true",
                   "--pids-limit", "256", "--tmpfs", "/tmp:size=16m,mode=1777",
                   "--publish", f"127.0.0.1:{port}:8096", "--mount", "type=volume,source=" + volume + ",target=/data",
                   "--env-file", str(env_file), image)
        wait_for_ready(request)
        assert docker("image", "inspect", "--format", "{{.Config.User}}", image) == "10001:10001", "Image must default to the non-root service user"
        assert docker("image", "inspect", "--format", "{{.Architecture}}", image) == platform.split("/")[1], "Image platform metadata is wrong"
        docker("exec", name, "/app/healthcheck")
        with request("/") as response:
            assert b'<div id="app">' in response.read(), "Static frontend missing from image"
        try:
            request("/api/apps")
            raise AssertionError("Anonymous protected read was accepted")
        except urllib.error.HTTPError as error:
            assert error.code == 401
        login()
        with request("/api/apps",{"name":"Container smoke"}) as response:
            app_id=json.load(response)["id"]
        with request("/api/apps/"+app_id+"/connection") as response:
            credentials=json.load(response)
        server=credentials["server"]
        with request("/api/apps/"+app_id+"/snapshot") as response:
            server_now = json.load(response)["now"]
        payload = smoke_payload(suffix, server_now)
        with request("/v1/logs", payload, server) as response:
            assert response.status == 200
        snapshot_has_record()
        # Reuse the maintained literal REST fixture generator/verifier against
        # this container's fresh database; successful receipt alone is not enough.
        fixture_env = dict(os.environ, SENTRAME_URL=origin, SENTRAME_ORIGIN=origin,
                           SENTRAME_PASSPHRASE=owner)
        repository = Path(__file__).resolve().parent.parent
        subprocess.run(["go", "run", "./cmd/fixtures", "standard"], cwd=repository,
                       env=fixture_env, check=True)
        manifest = json.loads((repository / "rest/generated/standard-manifest.json").read_text())
        with request("/api/apps") as response:
            demo_id=next(a["id"] for a in json.load(response) if a["name"]=="Sentrame demonstration")
        trace_path = "/api/apps/"+demo_id+"/traces/" + manifest["trace"]
        # The verifier has its own cookie jar; this owner cookie remains valid.
        with request(trace_path) as response:
            trace_before = json.load(response)
        trace_before.pop("revision", None)
        assert trace_before["mode"] == "trace" and trace_before["operations"]
        print("Fixture hashes:", json.dumps(manifest["hashes"], sort_keys=True))
        with request("/api/apps/"+app_id+"/events") as response:
            assert response.headers.get_content_type() == "text/event-stream"
            assert response.readline().startswith(b"event: reset")
        docker("restart", "--time", "30", name)
        wait_for_ready(request)
        jar.clear()
        login()
        snapshot_has_record()
        with request(trace_path) as response:
            trace_after = json.load(response)
        trace_after.pop("revision", None)
        assert trace_before == trace_after, "Native trace changed after persistent-volume restart"
        with request("/api/apps/"+app_id+"/connection") as response:
            assert json.load(response)==credentials,"Generated settings changed after restart"
        # A stopped coordinated archive preserves JSON, marker, every SQLite store
        # and numeric ownership. Round-trip the exact archive through the mount.
        docker("stop","--time","30",name)
        backup=Path(data_temporary.name)/"private-backup.tar"
        with backup.open("wb") as output:
            subprocess.run(["docker","cp","-a",name+":/data/.","-"],stdout=output,check=True)
        backup.chmod(0o600)
        broken=Path(data_temporary.name)/"broken-registry.json"
        broken.write_text("{");broken.chmod(0o600)
        docker("cp",str(broken),name+":/data/apps.json")
        docker("start",name)
        exit_code=subprocess.check_output(["docker","wait",name],text=True,timeout=20).strip()
        assert exit_code!="0","Damaged registry did not fail closed"
        with backup.open("rb") as source:
            subprocess.run(["docker","cp","-a","-",name+":/data"],stdin=source,check=True)
        docker("start",name);wait_for_ready(request);jar.clear();login();snapshot_has_record()
        with request("/api/apps/"+app_id+"/connection") as response:
            assert json.load(response)==credentials,"Coordinated restore changed authority"
        print(f"PASS {platform}: non-root read-only image, readiness, static UI, protected read, standard OTLP/HTTP JSON log, 37-step native-model REST verifier, protected trace relationships, SSE and exact persistent-volume trace evidence after restart; damaged-registry rejection and private coordinated backup/restore")
    except Exception:
        # This container contains only disposable synthetic data. Print bounded
        # startup diagnostics, never its environment or generated credentials.
        try:
            print("Container state:", docker("inspect", "--format", "{{.State.Status}} exit={{.State.ExitCode}}", name), flush=True)
            print(docker("logs", "--tail", "20", name), flush=True)
        except subprocess.CalledProcessError:
            pass
        raise
    finally:
        subprocess.run(["docker", "rm", "--force", name], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, check=False)
        subprocess.run(["docker","volume","rm",volume],stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL,check=False)
        data_temporary.cleanup()


if __name__ == "__main__":
    main()
