# OSMonitor

A small Go service that collects operating-system metrics and serves them as JSON over HTTP, plus a Nuxt UI dashboard in [ui/](ui/) that charts them live.

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
ui/                           Nuxt 4 + Nuxt UI 4 + Nuxt Charts dashboard
Dockerfile                    Go API image: multi-stage, static binary, non-root
docker-compose.yml            API + dashboard, host /proc and /sys mounted read-only
```

## Run locally

```sh
make run            # or: go run ./cmd/osmonitor
curl localhost:8080/api/metrics
```

On macOS you may see "You have not agreed to the Xcode and Apple SDKs license". Both `/usr/bin/make` and the cgo C compiler are Xcode shims that refuse to run while `xcode-select` points at Xcode.app. Point it at the Command Line Tools instead:

```sh
sudo xcode-select --switch /Library/Developer/CommandLineTools
```

Or skip `make` and run Go directly. The Makefile already disables cgo, so once the selector is switched `make run` works too:

```sh
CGO_ENABLED=0 go run ./cmd/osmonitor
```

## Endpoints

| Path                    | Description                                          |
|-------------------------|------------------------------------------------------|
| `GET /health`           | Liveness check                                       |
| `GET /api/metrics`      | Latest cached snapshot (refreshed every interval)    |
| `GET /api/metrics/live` | Collects a fresh snapshot on request (slower)        |

## Configuration

| Variable                            | Default | Description                                            |
|-------------------------------------|---------|--------------------------------------------------------|
| `OSMONITOR_ADDR`                    | `:8080` | Listen address                                         |
| `OSMONITOR_INTERVAL`                | `5s`    | Background refresh interval                            |
| `HOST_PROC`, `HOST_SYS`, `HOST_ETC` | unset   | Alternate procfs/sysfs paths (set in the Docker image) |

## Dashboard (ui/)

The `ui` folder is a Nuxt 4 app using [Nuxt UI](https://ui.nuxt.com) components and [Nuxt Charts](https://nuxtcharts.com). It polls the Go API every 3 seconds and shows CPU, memory and swap trends, per-core load, temperature sensors, disk usage and host details.

```sh
make ui-install     # npm install (first time)
make run            # Go API on :8080, in one terminal
make ui-dev         # Nuxt dev server on :3000, in another
```

The browser only ever calls its own origin. In development the Nuxt dev server proxies `/api/*` to the Go API; in production nginx does the same. Either way Go needs no CORS setup. Point the dev proxy elsewhere with `NUXT_API_URL`. If port 3000 is busy, run `npx nuxt dev --port 3210` inside `ui/`.

See [ui/README.md](ui/README.md) for the frontend layout.

## Docker

```sh
make docker-up      # docker compose up --build -d
open http://localhost:3000       # dashboard
curl localhost:8080/api/metrics  # raw JSON
make docker-down
```

Compose starts two containers: `osmonitor` (Go API, ~22 MB) and `ui` (static dashboard served by nginx, proxying `/api/*` to the API over the compose network). The UI upstream is set with `API_URL`. It mounts `/proc`, `/sys` and `/etc` from the host read-only and uses `pid: host`, so the API reports the host machine's metrics instead of its own cgroup view.

### Temperature caveats

- **Linux host**: sensors are read from `/sys/class/hwmon` and `/sys/class/thermal`. With the compose mounts this works inside the container.
- **macOS**: sensors are read directly when running the binary natively. Inside Docker Desktop the container sees the Linux VM, not the Mac, so `temperature` is empty.
- Readings outside -50..200 °C are discarded; some platforms return sentinel values for unpopulated sensor slots.

## Next steps (ideas)

- Add a `/metrics` endpoint in Prometheus format
- Add network and process metrics from gopsutil `net` and `process`
- Persist snapshots to a time-series store
