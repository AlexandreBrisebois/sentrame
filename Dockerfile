# syntax=docker/dockerfile:1
# Exact builder versions are verified against Docker Official Images. Resolve
# these tags to reviewed manifest digests before a reproducible published release.
ARG GO_IMAGE=golang:1.27.3-alpine3.24
ARG NODE_IMAGE=node:26.13.0-alpine3.24

FROM --platform=$BUILDPLATFORM ${NODE_IMAGE} AS frontend
WORKDIR /src/web
COPY web/package.json web/package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm ci
COPY web/ ./
RUN npm test && npm run build

FROM --platform=$BUILDPLATFORM ${GO_IMAGE} AS backend
ENV GOTOOLCHAIN=local CGO_ENABLED=0
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags='-s -w' -o /out/sentrame ./cmd/sentrame && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags='-s -w' -o /out/healthcheck ./cmd/healthcheck && \
    mkdir -p /out/data && chown 10001:10001 /out/data

FROM scratch AS runtime
WORKDIR /app
COPY --from=backend /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=backend /out/sentrame /app/sentrame
COPY --from=backend /out/healthcheck /app/healthcheck
COPY --from=backend --chown=10001:10001 /out/data /data
COPY --from=frontend /src/web/dist /app/web/dist
ENV SENTRAME_ADDR=0.0.0.0:8096 SENTRAME_DB=/data/sentrame.db
USER 10001:10001
EXPOSE 8096
HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=3 CMD ["/app/healthcheck"]
ENTRYPOINT ["/app/sentrame"]
