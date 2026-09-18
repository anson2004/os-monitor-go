# Multi-arch friendly: the Node and Go stages run on the build machine's own
# architecture ($BUILDPLATFORM) and Go cross-compiles for the target, so
# `docker buildx build --platform linux/amd64,linux/arm64` needs no emulation
# for the heavy steps. Only the tiny runtime stage is per-architecture.

# ---- dashboard build (architecture independent output) ----
FROM --platform=$BUILDPLATFORM node:22-alpine AS ui
WORKDIR /ui
COPY ui/package.json ui/package-lock.json ./
RUN npm ci
COPY ui/ .
RUN npm run generate

# ---- Go build (embeds the dashboard, cross-compiles) ----
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=ui /ui/.output/public ./internal/web/dist
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/osmonitor ./cmd/osmonitor

# ---- runtime (built per target platform) ----
FROM alpine:3.20
RUN adduser -D -u 10001 osmonitor
COPY --from=build /out/osmonitor /usr/local/bin/osmonitor

# gopsutil reads the host's procfs/sysfs from these paths when mounted.
ENV HOST_PROC=/host/proc \
    HOST_SYS=/host/sys \
    HOST_ETC=/host/etc \
    OSMONITOR_ADDR=:8080 \
    OSMONITOR_INTERVAL=5s

USER osmonitor
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD wget -qO- http://localhost:8080/health || exit 1
ENTRYPOINT ["osmonitor"]
