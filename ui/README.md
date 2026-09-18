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
server/api/[...path].ts        forwards /api/* to the Go API (no CORS needed)
Dockerfile                     multi-stage build, runs the Nitro server on :3000
```

## Develop

Start the Go API first (`make run` in the parent folder), then:

```sh
npm install
npm run dev            # http://localhost:3000
```

If port 3000 is taken: `npx nuxt dev --port 3210`.

The Go API location defaults to `http://localhost:8080` and can be changed with
`NUXT_API_URL`.

## Check

```sh
npm run lint
npm run typecheck
```

## Build

```sh
npm run build          # output in .output/
node .output/server/index.mjs
```

Or via Docker from the parent folder: `docker compose up --build`.
