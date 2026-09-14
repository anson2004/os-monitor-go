# OSMonitor

A small Go service that collects operating-system metrics and serves them as JSON over HTTP.

Collected:

- **Host**: hostname, OS, kernel, architecture, uptime, process count, load average
- **CPU**: model, core counts, frequency, total and per-core usage
- **Memory**: RAM and swap usage
- **Disks**: usage per mounted filesystem
- **Temperature**: hardware sensor readings (platform dependent, see below)

Metrics come from [gopsutil](https://github.com/shirou/gopsutil), so the same binary runs on Linux, macOS and Windows.

## Layout

```
cmd/osmonitor/main.go         entrypoint: config from env, graceful shutdown
internal/collector/           one file per metric group, combined into a Snapshot
internal/server/server.go     background refresh loop + HTTP handlers
Dockerfile                    multi-stage build, static binary, non-root user
docker-compose.yml            runs with host /proc and /sys mounted read-only
```

## Run locally

```sh
make run            # or: go run ./cmd/osmonitor
curl localhost:8080/api/metrics
```

On macOS you may see "You have not agreed to the Xcode and Apple SDKs license". Both `/usr/bin/make` and the cgo C compiler are Xcode shims that refuse to run until the license is accepted. Either accept it once:

```sh
sudo xcodebuild -license accept
```

or skip `make` and run Go directly (the Makefile already disables cgo, so after accepting the license `make run` works too):

```sh
CGO_ENABLED=0 go run ./cmd/osmonitor
```

## Endpoints

| Path                 | Description                                              |
|----------------------|----------------------------------------------------------|
| `GET /health`        | Liveness check                                           |
| `GET /api/metrics`   | Latest cached snapshot (refreshed every interval)        |
| `GET /api/metrics/live` | Collects a fresh snapshot on request (slower)         |

## Configuration

| Variable              | Default | Description                       |
|-----------------------|---------|-----------------------------------|
| `OSMONITOR_ADDR`      | `:8080` | Listen address                    |
| `OSMONITOR_INTERVAL`  | `5s`    | Background refresh interval       |
| `HOST_PROC`, `HOST_SYS`, `HOST_ETC` | unset | Alternate procfs/sysfs paths (set in the Docker image) |

## Docker

```sh
make docker-up      # docker compose up --build -d
curl localhost:8080/api/metrics
make docker-down
```

The compose file mounts `/proc`, `/sys` and `/etc` from the host read-only and uses `pid: host`, so the container reports the host machine's metrics instead of its own cgroup view.

### Temperature caveats

- **Linux host**: sensors are read from `/sys/class/hwmon` and `/sys/class/thermal`. With the compose mounts this works inside the container.
- **macOS**: sensors are read directly when running the binary natively. Inside Docker Desktop the container sees the Linux VM, not the Mac, so `temperature` is empty.
- Readings outside -50..200 °C are discarded; some platforms return sentinel values for unpopulated sensor slots.

## Next steps (ideas)

- Add a `/metrics` endpoint in Prometheus format
- Add network and process metrics from gopsutil `net` and `process`
- Persist snapshots to a time-series store
