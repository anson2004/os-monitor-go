# ---- build stage ----
FROM golang:1.24-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/osmonitor ./cmd/osmonitor

# ---- runtime stage ----
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
