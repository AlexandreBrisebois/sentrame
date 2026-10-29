# Multi-architecture container and NAS launch

The container package targets **linux/amd64 and linux/arm64** under the same Docker Hub tag: `brisebois/sentrame`. Docker selects the matching platform when the NAS pulls it. The browser assets, Go service, maintained OTLP receiver and SQLite run in one non-root container; there is no separate frontend server or public database port.

This package is prepared locally. A successful GitHub build/smoke/publish run is required before claiming either container platform qualified. It does not establish Home Server capacity, Cloudflare Tunnel behavior, browser SDK compatibility, power-loss durability or a WFS integration.

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

Copy `compose.yaml` and `.env.example` into the chosen local NAS project directory, then:

```sh
cp .env.example .env
chmod 600 .env
```

Edit `.env` locally:

- Set `SENTRAME_IMAGE` to the published tag or, preferably, index digest.
- Set `SENTRAME_ORIGIN` to the exact URL used by browsers, such as `https://sentrame.your-domain.example`. It controls same-origin mutation checks and secure cookies; do not use a wildcard or a URL path.
- Set a fresh owner passphrase of at least 12 characters and two distinct random producer tokens. Use `openssl rand -hex 32` locally for each token. Never reuse the synthetic Mac credentials for deployment.
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

Replace the example with the folder actually selected for this deployment. Do not run these commands against an unrelated or shared folder. Compose requires the bind folder to exist; it does not create a root-owned replacement silently. The folder holds `sentrame.db`, its WAL, and the process lease across container replacement.

Then:

```sh
docker compose config --quiet
docker compose pull
docker compose up -d
docker compose ps
docker compose exec sentrame /app/healthcheck
```

The root filesystem is read-only; `/tmp` is a bounded tmpfs. The process runs as UID/GID 10001 with dropped capabilities and a 256-PID limit. No RAM/CPU reservation or NAS fit is inferred. Retention is 30-day detail and 12-month supported summaries; capacity/disk monitoring and backups are still required.

Keep a coherent SQLite backup. Stopping Sentrame before copying the whole selected data folder is the simplest backup/rollback boundary. A running database must use a SQLite-aware backup including its WAL state, not an arbitrary database-file copy.

## Connect the existing tunnel

The default host port is **127.0.0.1:8096**, so it is not exposed to the LAN. There is no tunnel container/token/DNS setup in this Compose project.

- If the existing connector runs on the NAS host, its origin can be `http://127.0.0.1:8096`.
- If the connector runs in Docker, its `localhost` is a different container. Deliberately attach Sentrame and the connector to the existing appropriate Docker network, then use `http://sentrame:8096` as the origin. Review the actual network/connector configuration first; this package does not auto-attach to a shared network or create another tunnel.
- A reverse proxy must preserve the browser's `Origin` header. The configured `SENTRAME_ORIGIN` must match the public HTTPS URL. Sentrame does not trust forwarded headers for authorization or app attribution.

Verify the real HTTPS hostname, secure cookies, protected reads, CSRF rejection, live SSE flushing/reconnect and temporary-viewer expiry/revocation through that route before sharing access. HTTPS reachability alone does not qualify the route. Public-hostname Tunnel gRPC is not advertised; use the separately verified OTLP/HTTP paths.

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

Keep the previous image digest and stopped-data backup. SQLite schema additions are not a general downgrade guarantee: roll back the corresponding database snapshot with the prior image when compatibility is unknown. Restoring old backups can restore deleted evidence or old grant metadata; treat backups as a separate retention/access boundary. Viewer sessions are in memory and end on service restart; already redeemed one-use links remain used, so an owner issues a fresh invite after restart. The selected bind folder remains on disk when the container is replaced or removed; deleting that folder or its backup is a separate destructive action.

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

Local preparation passed Compose configuration validation, YAML syntax/duplicate-key checks, Python probe regression tests and script compilation. The service and health probe cross-compile to static Linux amd64/arm64 ELF binaries. The sandbox denied Docker daemon access. The reviewer ran `task test && task container:smoke`: the test task passed and the ARM64 image built successfully. Two initial smoke checks exposed test-harness defects: startup connection resets were not retried, and a host-clock fixture could fall outside the container’s occurrence-time window. Both were fixed with regression coverage, retaining the accepted-unassigned assertion. The reviewer then supplied successful `PASS linux/arm64` and `PASS linux/amd64` results: non-root/read-only startup, private bind-folder storage, readiness/static UI, protected reads, standard OTLP/HTTP JSON logs, SSE and persistence after restart all passed on Docker Desktop. This is local runtime evidence (including AMD64 execution on the ARM64 Mac), not native NAS workload or routed-transport qualification. No workflow has been dispatched, Docker Hub image published or NAS deployed.

Official source checks on 2026-10-28 verified `golang:1.27.3-alpine3.24` and `node:26.13.0-alpine3.24` in Docker's maintained image manifest lists with amd64/arm64 support. Registry digest resolution was unavailable locally, so the workflow resolves version tags once and uses those same immutable builder references for both smoke builds and final publication. Local frontend tests use installed Node 26.12.0; the workflow independently tests/builds on the pinned 26.13.0 builder version.

- [Docker multi-platform GitHub Actions guidance](https://docs.docker.com/build/ci/github-actions/multi-platform/).
- [Go Official Image source manifest](https://raw.githubusercontent.com/docker-library/official-images/master/library/golang) and [Node Official Image source manifest](https://raw.githubusercontent.com/docker-library/official-images/master/library/node).
- [Registry manifest inspection](https://docs.docker.com/reference/cli/docker/buildx/imagetools/inspect/).
- Workflow action commits were resolved from their official release pages: [checkout7.0.1](https://github.com/actions/checkout/releases/tag/v7.0.1), [setup-go7.0.0](https://github.com/actions/setup-go/releases/tag/v7.0.0), [setup-node7.1.0](https://github.com/actions/setup-node/releases/tag/v7.1.0), [build-push7.4.0](https://github.com/docker/build-push-action/releases/tag/v7.4.0), [buildx4.4.1](https://github.com/docker/setup-buildx-action/releases/tag/v4.4.1), [QEMU4.4.0](https://github.com/docker/setup-qemu-action/releases/tag/v4.4.0), and [login4.7.0](https://github.com/docker/login-action/releases/tag/v4.7.0).

- [GitHub manual workflow discovery and default-branch requirement](https://docs.github.com/en/actions/how-tos/manage-workflow-runs/manually-run-a-workflow).

The owner has authorized committing/pushing the implementation branch after the successful local smoke checks. Configure the DockerHub secret before creating a publication tag; a branch push alone does not publish. No release tag is created as part of this preparation.
