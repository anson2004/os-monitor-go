<script setup lang="ts">
import type { Temperature } from '#shared/types/metrics'

const props = defineProps<{ sensors: Temperature[] }>()

const sorted = computed(() =>
  [...props.sensors].sort((a, b) => b.celsius - a.celsius)
)

const hottest = computed(() => sorted.value[0])

function tempColor(c: number): 'success' | 'warning' | 'error' {
  if (c >= 85) return 'error'
  if (c >= 70) return 'warning'
  return 'success'
}
</script>

<template>
  <UEmpty
    v-if="!sensors.length"
    icon="i-lucide-thermometer-snowflake"
    title="No sensors available"
    description="Temperature sensors are not exposed on this host. Inside Docker this requires mounting the host /sys."
  />

  <div
    v-else
    class="space-y-4"
  >
    <div
      v-if="hottest"
      class="flex items-baseline gap-2"
    >
      <span class="text-3xl font-semibold tabular-nums">{{ hottest.celsius.toFixed(1) }}°C</span>
      <span class="text-sm text-muted">hottest · {{ hottest.sensor }}</span>
    </div>

    <ul class="max-h-64 space-y-2 overflow-y-auto pr-1">
      <li
        v-for="s in sorted"
        :key="s.sensor"
        class="flex items-center gap-3 text-sm"
      >
        <span
          class="w-32 truncate text-muted"
          :title="s.sensor"
        >{{ s.sensor }}</span>
        <UProgress
          :model-value="s.celsius"
          :max="s.critical || 100"
          :color="tempColor(s.celsius)"
          size="sm"
          class="flex-1"
        />
        <span class="w-16 text-right tabular-nums">{{ s.celsius.toFixed(1) }}°C</span>
      </li>
    </ul>
  </div>
</template>
