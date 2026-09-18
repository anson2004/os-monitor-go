# OSMonitor UI

Dashboard for the Go metrics API in the parent folder. Built with
[Nuxt 4](https://nuxt.com), [Nuxt UI 4](https://ui.nuxt.com) and
[Nuxt Charts 3](https://nuxtcharts.com).

## Layout

```
app/app.vue                    shell: header, colour-mode toggle
app/pages/index.vue            dashboard grid
app/components/                StatCard, UsageHistoryChart, PerCoreChart,
                               TemperaturePanel, DiskTable, HostCard
app/composables/useMetrics.ts  polls /api/metrics, keeps a 60-point history
app/utils/format.ts            bytes / percent / uptime formatting
shared/types/metrics.ts        TypeScript mirror of the Go Snapshot JSON
nginx/default.conf.template   nginx: serves the static build, proxies /api/* to Go
Dockerfile                     node build stage -> nginx:alpine runtime on :80
```

## Develop

Start the Go API first (`make run` in the parent folder), then:

```sh
npm install
npm run dev            # http://localhost:3000
```

If port 3000 is taken: `npx nuxt dev --port 3210`.

In development the Nuxt dev server proxies `/api/*` to the Go API (`nitro.devProxy` in `nuxt.config.ts`), so the browser stays same-origin. Default target is `http://localhost:8080`; override with `NUXT_API_URL`.

## Check

```sh
npm run lint
npm run typecheck
```

## Build

The app is a static single-page build (`ssr: false`):

```sh
npm run generate       # static site in .output/public/
```

In production nginx serves those files and proxies `/api/*` to the Go API. The upstream is set with `API_URL` (default `http://osmonitor:8080`), substituted into `nginx/default.conf.template` when the container starts.

```sh
docker build -t osmonitor-ui .
docker run -p 3000:80 -e API_URL=http://host.docker.internal:8080 osmonitor-ui
```

Or from the parent folder: `docker compose up --build`.
