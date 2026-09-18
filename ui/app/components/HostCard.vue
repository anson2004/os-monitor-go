<script setup lang="ts">
import type { HostInfo } from '#shared/types/metrics'

const props = defineProps<{ host: HostInfo }>()

const fields = computed(() => [
  { label: 'Hostname', value: props.host.hostname },
  { label: 'OS', value: `${props.host.platform || props.host.os} ${props.host.platform_version}`.trim() },
  { label: 'Kernel', value: props.host.kernel_version },
  { label: 'Architecture', value: props.host.arch },
  { label: 'Processes', value: String(props.host.processes) },
  { label: 'Load (1 / 5 / 15 min)', value: `${props.host.load_1} / ${props.host.load_5} / ${props.host.load_15}` }
])
</script>

<template>
  <UCard>
    <template #header>
      <h2 class="font-medium">
        Host
      </h2>
    </template>
    <dl class="grid gap-x-6 gap-y-3 text-sm sm:grid-cols-2 lg:grid-cols-3">
      <div
        v-for="f in fields"
        :key="f.label"
      >
        <dt class="text-muted">
          {{ f.label }}
        </dt>
        <dd class="font-medium break-all">
          {{ f.value || '—' }}
        </dd>
      </div>
    </dl>
  </UCard>
</template>
