# OSMonitor UI

Dashboard for the Go metrics API in the parent folder. Built with
[Nuxt 4](https://nuxt.com), [Nuxt UI 4](https://ui.nuxt.com) and
[Nuxt Charts 3](https://nuxtcharts.com). It is a static single-page app
(`ssr: false`) that the Go binary embeds and serves.

## Layout

```
app/app.vue                    shell: header, colour-mode toggle
app/pages/index.vue            dashboard grid
app/components/                StatCard, UsageHistoryChart, PerCoreChart,
                               TemperaturePanel, DiskTable, HostCard
app/composables/useMetrics.ts  polls /api/metrics, keeps a 60-point history
app/utils/format.ts            bytes / percent / uptime formatting
shared/types/metrics.ts        TypeScript mirror of the Go Snapshot JSON
```

## Develop

Start the Go API first (`make run` in the parent folder), then:

```sh
npm install
npm run dev            # http://localhost:3000
```

If port 3000 is taken: `npx nuxt dev --port 3210`.

The dev server proxies `/api/*` to the Go API (`nitro.devProxy` in
`nuxt.config.ts`), so the browser stays same-origin. Default target is
`http://localhost:8080`; override with `NUXT_API_URL`.

## Check

```sh
npm run lint
npm run typecheck
```

## Build

```sh
npm run generate       # static site in .output/public/
```

From the parent folder, `make ui-build` runs this and copies the output into
`internal/web/dist/`, where the Go build embeds it. The root `Dockerfile` does
the same in its first stage, so `docker compose up --build` needs no manual step.
