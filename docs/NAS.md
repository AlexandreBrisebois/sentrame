# Multi-architecture container and NAS launch

Current multi-app source uses greenfield registration in a fresh directory. Follow [App workspaces](MULTI-APP.md) for exact registry persistence, permissions, pending-removal recovery, coordinated backup/restore and rollback limitations. Existing singleton files are never imported. Historical published images listed below do not include this implementation.

The container package targets **linux/amd64 and linux/arm64** under the same Docker Hub tag: `brisebois/sentrame`. Docker selects the matching platform when the NAS pulls it. The browser assets, Go service, maintained OTLP receiver and SQLite run in one non-root container; there is no separate frontend server or public database port.

Release candidate `brisebois/sentrame:v0.1.0-rc.1` is published. Its [GitHub workflow](https://github.com/AlexandreBrisebois/sentrame/actions/runs/40808148393) passed checks, both platform smoke jobs, publication and the final manifest-platform check. It does not establish Home Server capacity, Cloudflare Tunnel behavior, browser SDK compatibility, power-loss durability or a WFS integration.

## Publish through GitHub Actions

Configure the `AlexandreBrisebois/sentrame` repository (or its `dockerhub` environment):

| Setting | Value |
| --- | --- |
| Actions variable `DOCKERHUB_USERNAME` | Docker Hub login allowed to push `brisebois/sentrame`; defaults to `brisebois` |
| Actions secret `DOCKERHUB_TOKEN` | Docker Hub access token with write access to this repository |
| Docker Hub repository | `brisebois/sentrame`, with the visibility you intend |

Do not put the Docker token, owner passphrase or producer credentials in workflow files or build arguments. An environment named `dockerhub` scopes the publishing job; protection/approval rules are optional settings controlled by the repository owner.

The workflow has three gates:

1. Run Go behavior tests/vet and frontend tests/type checking/production build. Resolve exact official builder version tags to immutable manifest digests once for the run.
2. Build and run each platform separately. The arm64 runtime is exercised under QEMU on the amd64 runner. The disposable smoke test checks the non-root/read-only image, readiness, static app, protected API, standard OTLP/HTTP JSON log receipt, SSE and the same committed evidence after restart on a local bind folder.
3. Publish one multi-platform index only after both platform smoke jobs pass. Include build provenance and an SBOM, inspect the resulting index, and record its digest in the job summary.

Pull requests run checks and smoke tests **without registry login or publication**. A pushed `v*` tag publishes that exact tag and `sha-<12-character commit>`. There is no automatic `latest`. Once the workflow exists on the repository’s default branch, alternatively run **Actions → Container → Run workflow**, choose `publish=true` and an explicit review/release tag. The manual Run workflow UI requires default-branch workflow discovery; the version-tag push path can run the committed workflow without merging this implementation branch. A manual run defaults to build/test only. GitHub can execute only workflow/code that has been committed and pushed; this local preparation does not commit, push, create a release/tag or run the workflow.

Use the emitted index digest for deployment, for example:

```dotenv
SENTRAME_IMAGE=brisebois/sentrame@sha256:REPLACE_WITH_THE_PUBLISHED_INDEX_DIGEST
```

Version tags may be overwritten by a later authorized run; a digest identifies the selected artifact. Record the previous digest and database backup before upgrading.

## Prepare the NAS

The NAS needs Docker Engine/Compose with either an amd64 or arm64 Linux host and a **local** filesystem for SQLite. Do not put `/data` on an SMB/NFS mount. No NAS hardware, filesystem, resource availability or existing tunnel has been inspected here.

Choose one standalone Compose example:

| File | Access | Start command |
| --- | --- | --- |
| `compose.yaml` | Host loopback port8096, configurable with `SENTRAME_PORT` | `docker compose up -d` |
| `compose.cloudflare.yaml` | Existing external `cloudflare-proxy` network, no host port | `docker compose -f compose.cloudflare.yaml up -d` |

Copy the chosen file and `.env.example` into your deployment directory. Use one Compose file per deployment; the Cloudflare example is a complete configuration, not an override layered onto the standard file.

Then:

```sh
cp .env.example .env
chmod 600 .env
```

Edit `.env` locally:

- Set `SENTRAME_IMAGE` to the published tag or, preferably, index digest.
- Keep `.env` alongside the selected Compose file; Compose uses it for substitutions and passes its application settings into the container. Set the owner passphrase before starting. Register apps afterward in the UI; configuration syntax validation does not validate credential strength.
- For the standard example, use `SENTRAME_ORIGIN=http://127.0.0.1:8096` (adjust its port if changing `SENTRAME_PORT`). For the Cloudflare example, set `SENTRAME_ORIGIN` to your exact public HTTPS origin, such as `https://sentrame.your-domain.example`. It controls same-origin mutation checks and secure cookies; do not use a wildcard or a URL path.
- Set a fresh owner passphrase of at least 12 characters. Register apps through the UI to generate producer credentials; do not set producer token environment variables.
- Keep these secrets on the NAS and in authorized exporter configuration. Container environment variables are visible to administrators with Docker access.

Compose binds the **specific local folder** in `SENTRAME_DATA_DIR` to `/data`. For example, choose a new dedicated directory `/srv/sentrame/data` and set:

```dotenv
SENTRAME_DATA_DIR=/srv/sentrame/data
```

Before starting, have the NAS administrator prepare only that chosen directory:

```sh
sudo mkdir -p /srv/sentrame/data
sudo chown 10001:10001 /srv/sentrame/data
sudo chmod 700 /srv/sentrame/data
```

Replace the example with the folder actually selected for this deployment. Do not run these commands against an unrelated or shared folder. Compose requires the bind folder to exist; it does not create a root-owned replacement silently. The folder holds private `apps.json`, `.apps-initialized`, app stores under `apps/<id>/telemetry.db`, WAL/SHM files and leases across container replacement.

Then:

```sh
docker compose config --quiet
docker compose pull
docker compose up -d
docker compose ps
docker compose exec sentrame /app/healthcheck
```

Both minimal Compose examples inherit the image’s non-root UID/GID10001, listening address, database path and health check. It loads the application settings from `.env` and adds the selected data bind and existing proxy network. Optional runtime isolation/resource settings are left to the deployment environment; no capability-drop setting is required by this example. Retention is30-day detail and12-month supported summaries; capacity/disk monitoring and backups are still required.

Keep a coherent SQLite backup. Stopping Sentrame before copying the whole selected data folder is the simplest backup/rollback boundary. A running database must use a SQLite-aware backup including its WAL state, not an arbitrary database-file copy.

## Connect the existing proxy network

`compose.cloudflare.yaml` attaches Sentrame to the existing external Docker network **cloudflare-proxy**. That network and the Cloudflare connector/proxy must already exist. This project does not create a second connector, manage its token or configure DNS.

In the existing proxy, point the Sentrame hostname at:

```text
http://sentrame:8096
```

Start that configuration with `docker compose -f compose.cloudflare.yaml up -d`. Include `-f compose.cloudflare.yaml` in its `config`, `pull`, `ps`, `exec` and `stop` commands as well.

Both the proxy and Sentrame must be attached to `cloudflare-proxy` for that container name to resolve. There is no host `ports` mapping; use the public hostname through the proxy to open Sentrame. The origin link above is HTTP within the shared network; the browser-facing hostname should use HTTPS.

Set `SENTRAME_ORIGIN` in `.env` to the exact public browser origin, for example `https://sentrame.example.com`, without a trailing slash or URL path. The proxy must preserve the browser's `Origin` header. Sentrame does not use forwarded headers to grant producer or viewer authority.

Check the real hostname, secure cookies, protected reads, CSRF rejection, live SSE flushing/reconnect and temporary-viewer expiry/revocation through the route before sharing access. HTTPS reachability alone does not qualify the route. Public-hostname Tunnel gRPC is not advertised; use the separately verified OTLP/HTTP paths.

## Prepare OTel senders after route verification

The published app serves `/v1/logs`, `/v1/traces` and `/v1/metrics` over OTLP/HTTP. Configure an ordinary standard exporter with the base endpoint and `Authorization: Bearer <producer-token>` header. The server token is for trusted server producers; the browser token cannot establish server authority. Never embed the server token or viewer passphrase in a browser bundle.

Start with a synthetic standard request through the actual NAS route and inspect accepted-unassigned receipt before adding business-event mappings. Real .NET/Next.js/browser exporters, CORS/CSP/browser lifecycle, and WFS event/privacy choices each need their own qualification. The container does not install instrumentation or contact WFS. This repository's synthetic REST/fixture demonstrations remain the repeatable reference; run them against a deliberately selected test destination and credentials, not real app traffic by accident.

## Upgrade and rollback

```sh
docker compose stop
# Back up the complete selected data folder while stopped, then set the new image digest.
docker compose pull
docker compose up -d
docker compose exec sentrame /app/healthcheck
```

Keep the previous image digest and stopped-data backup. The greenfield structural model requires SQLite schema3 / structural-v2. Earlier schemas are rejected before initialization; no legacy reader, migration or automatic reset is provided. Preserve the old directory and select a fresh dedicated test directory for cutover. Rollback uses the corresponding stopped database snapshot and prior image, never a newer database under an older image. Restoring old backups can restore deleted evidence or old grant metadata; treat backups as a separate retention/access boundary. Viewer sessions are in memory and end on service restart; already redeemed one-use links remain used, so an owner issues a fresh invite after restart. The selected bind folder remains on disk when the container is replaced or removed; deleting that folder or its backup is a separate destructive action.

## Run a local container check

When Docker is available on the Mac, build and smoke only the desired local platform:

```sh
docker buildx build --platform linux/arm64 --load -t sentrame:review .
IMAGE_UNDER_TEST=sentrame:review PLATFORM_UNDER_TEST=linux/arm64 python3 scripts/container-smoke.py
```

For a local multi-platform export without publishing:

```sh
docker buildx build --platform linux/amd64,linux/arm64 --output type=oci,dest=/private/tmp/sentrame-review.oci.tar .
```

The smoke script uses fresh synthetic credentials and its own uniquely named container and temporary bind folder; cleanup removes only those test resources. The test folder alone permits UID 10001 writes through disposable permissions under a private parent; deployment still uses the chosen UID/GID 10001-owned mode 700 folder. It does not read or reset the Mac review database. The default Docker image-store cannot necessarily load a multi-platform index; use the OCI export or the workflow rather than claiming a two-platform `--load` succeeded.

## Verification and sources

Local preparation passed Compose configuration validation, YAML syntax/duplicate-key checks, Python probe regression tests and script compilation. The service and health probe cross-compile to static Linux amd64/arm64 ELF binaries. The sandbox denied Docker daemon access. The reviewer ran `task test && task container:smoke`: the test task passed and the ARM64 image built successfully. Two initial smoke checks exposed test-harness defects: startup connection resets were not retried, and a host-clock fixture could fall outside the container’s occurrence-time window. Both were fixed with regression coverage, retaining the accepted-unassigned assertion. The reviewer then supplied successful `PASS linux/arm64` and `PASS linux/amd64` results: non-root/read-only startup, private bind-folder storage, readiness/static UI, protected reads, standard OTLP/HTTP JSON logs, SSE and persistence after restart all passed on Docker Desktop. This is local runtime evidence (including AMD64 execution on the ARM64 Mac), not native NAS workload or routed-transport qualification. The initial GitHub workflow has since passed all jobs and published the release candidate; the NAS itself has not been deployed or qualified.

Official source checks on 2026-10-28 verified `golang:1.27.3-alpine3.24` and `node:26.13.0-alpine3.24` in Docker's maintained image manifest lists with amd64/arm64 support. Registry digest resolution was unavailable locally, so the workflow resolves version tags once and uses those same immutable builder references for both smoke builds and final publication. Local frontend tests use installed Node 26.12.0; the workflow independently tests/builds on the pinned 26.13.0 builder version.

- [Docker multi-platform GitHub Actions guidance](https://docs.docker.com/build/ci/github-actions/multi-platform/).
- [Go Official Image source manifest](https://raw.githubusercontent.com/docker-library/official-images/master/library/golang) and [Node Official Image source manifest](https://raw.githubusercontent.com/docker-library/official-images/master/library/node).
- [Registry manifest inspection](https://docs.docker.com/reference/cli/docker/buildx/imagetools/inspect/).
- Workflow action commits were resolved from their official release pages: [checkout7.0.1](https://github.com/actions/checkout/releases/tag/v7.0.1), [setup-go7.0.0](https://github.com/actions/setup-go/releases/tag/v7.0.0), [setup-node7.1.0](https://github.com/actions/setup-node/releases/tag/v7.1.0), [build-push7.4.0](https://github.com/docker/build-push-action/releases/tag/v7.4.0), [buildx4.4.1](https://github.com/docker/setup-buildx-action/releases/tag/v4.4.1), [QEMU4.4.0](https://github.com/docker/setup-qemu-action/releases/tag/v4.4.0), and [login4.7.0](https://github.com/docker/login-action/releases/tag/v4.7.0).

- [GitHub manual workflow discovery and default-branch requirement](https://docs.github.com/en/actions/how-tos/manage-workflow-runs/manually-run-a-workflow).

The initial implementation was committed as `5738712b7469d5ef24810f53a2af1d71c9961925`, pushed, and tagged `v0.1.0-rc.1`. Publication succeeded after registry-secret configuration. Future branch pushes alone do not publish; use the workflow’s release-tag or manual-publication path.

### Published first release candidate

- Image: `brisebois/sentrame:v0.1.0-rc.1`
- Immutable index: `brisebois/sentrame@sha256:4524464d90d9d72e54ffc17b9ee8197a29380673cb1721fce3452a69d974891e`
- Platforms verified by the publishing job: `linux/amd64`, `linux/arm64`.
- [Successful workflow and logs](https://github.com/AlexandreBrisebois/sentrame/actions/runs/40808148393).

This records container and registry evidence, not a native NAS capacity test or a deployed external route. The environment example now selects this published release tag; use the index digest above when you want an immutable deployment selection.

## Structural model capability boundary

The published `v0.1.0-rc.1` predates standard trace exploration and is not a minimum supporting image for this change. The current local source uses schema3 / structural-v2 and journey-native-v3; no new image has been published by this implementation. Select a separately reviewed future image/build containing those identities before expecting the new trace/mapping/hint behavior. Local model qualification is recorded per architecture in [the evidence ledger](verification/standard-otel.md); previous container evidence above remains historical.

A sender on the same Mac uses its loopback origin. A sender in another container uses an explicitly reachable service address such as `http://sentrame:8096`; its own `127.0.0.1` is not Sentrame. An HTTP exporter base has no signal path; a signal-specific URL uses `/v1/traces` or `/v1/logs`. Configure `http/protobuf` explicitly, never infer gRPC from an open port. HTTP JSON/gzip and metrics receipt claims are limited to their individual tested cells. Browser exporters use reviewed bounded app-origin forwarding under browser authority. Direct Mac exporter passes do not qualify the NAS, shared Tunnel, Collector or HTTPS route. Qualify each actual route before collection, using fresh disposable data and private credentials.
