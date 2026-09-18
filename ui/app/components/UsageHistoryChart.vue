<script setup lang="ts">
import type { HistoryPoint } from '#shared/types/metrics'
import { formatTime } from '~/utils/format'

const props = defineProps<{ history: HistoryPoint[] }>()

const categories: Record<string, BulletLegendItemInterface> = {
  cpu: { name: 'CPU %', color: '#3b82f6' },
  memory: { name: 'Memory %', color: '#22c55e' },
  swap: { name: 'Swap %', color: '#f59e0b' }
}

const xFormatter = (i: number) => {
  const point = props.history[i]
  return point ? formatTime(point.time) : ''
}
const yFormatter = (v: number) => `${v}%`
</script>

<template>
  <div
    v-if="history.length < 2"
    class="flex h-[260px] items-center justify-center text-sm text-muted"
  >
    Collecting samples…
  </div>
  <AreaChart
    v-else
    :data="history"
    :categories="categories"
    :height="260"
    :y-domain="[0, 100]"
    :y-num-ticks="5"
    :x-num-ticks="6"
    :x-formatter="xFormatter"
    :y-formatter="yFormatter"
    :curve-type="CurveType.MonotoneX"
    :legend-position="LegendPosition.TopRight"
    :y-grid-line="true"
    :gradient="true"
  />
</template>
