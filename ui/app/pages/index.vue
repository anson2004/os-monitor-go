<script setup lang="ts">
import { formatBytes, formatPercent, formatTime, formatUptime } from '~/utils/format'

const { snapshot, history, error, loading, lastUpdated } = useMetrics(3000)
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-wrap items-center justify-between gap-2">
      <div>
        <h1 class="text-2xl font-semibold">
          {{ snapshot?.host.hostname ?? 'Dashboard' }}
        </h1>
        <p class="text-sm text-muted">
          <template v-if="snapshot">
            {{ snapshot.host.platform || snapshot.host.os }} {{ snapshot.host.platform_version }}
            · {{ snapshot.host.arch }}
          </template>
          <template v-else>
            Waiting for first sample…
          </template>
        </p>
      </div>
      <UBadge
        v-if="lastUpdated"
        color="neutral"
        variant="subtle"
        icon="i-lucide-refresh-cw"
      >
        Updated {{ formatTime(lastUpdated.toISOString()) }}
      </UBadge>
    </div>

    <UAlert
      v-if="error"
      color="error"
      variant="subtle"
      icon="i-lucide-triangle-alert"
      title="Cannot reach the metrics API"
      :description="error"
    />

    <UAlert
      v-if="snapshot?.errors?.length"
      color="warning"
      variant="subtle"
      icon="i-lucide-info"
      title="Some metrics are unavailable"
      :description="snapshot.errors.join(' · ')"
    />

    <div
      v-if="loading && !snapshot"
      class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4"
    >
      <USkeleton
        v-for="i in 4"
        :key="i"
        class="h-28"
      />
    </div>

    <template v-if="snapshot">
      <!-- Headline numbers -->
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard
          title="CPU"
          icon="i-lucide-cpu"
          :value="formatPercent(snapshot.cpu.usage_percent)"
          :percent="snapshot.cpu.usage_percent"
          :subtitle="`${snapshot.cpu.logical_cores} cores · ${snapshot.cpu.model_name || 'unknown'}`"
        />
        <StatCard
          title="Memory"
          icon="i-lucide-memory-stick"
          :value="formatPercent(snapshot.memory.used_percent)"
          :percent="snapshot.memory.used_percent"
          :subtitle="`${formatBytes(snapshot.memory.used_bytes)} of ${formatBytes(snapshot.memory.total_bytes)}`"
        />
        <StatCard
          title="Swap"
          icon="i-lucide-hard-drive-download"
          :value="formatPercent(snapshot.memory.swap_used_percent)"
          :percent="snapshot.memory.swap_used_percent"
          :subtitle="`${formatBytes(snapshot.memory.swap_used_bytes)} of ${formatBytes(snapshot.memory.swap_total_bytes)}`"
        />
        <StatCard
          title="Uptime"
          icon="i-lucide-timer"
          :value="formatUptime(snapshot.host.uptime_seconds)"
          :subtitle="`${snapshot.host.processes} processes · load ${snapshot.host.load_1} / ${snapshot.host.load_5} / ${snapshot.host.load_15}`"
        />
      </div>

      <!-- Trends and per-core -->
      <div class="grid gap-4 lg:grid-cols-3">
        <UCard class="lg:col-span-2">
          <template #header>
            <h2 class="font-medium">
              Usage over time
            </h2>
          </template>
          <UsageHistoryChart :history="history" />
        </UCard>

        <UCard>
          <template #header>
            <h2 class="font-medium">
              Per-core load
            </h2>
          </template>
          <PerCoreChart :per-core="snapshot.cpu.per_core_percent" />
        </UCard>
      </div>

      <!-- Temperature and disks -->
      <div class="grid gap-4 lg:grid-cols-2">
        <UCard>
          <template #header>
            <h2 class="font-medium">
              Temperature
            </h2>
          </template>
          <TemperaturePanel :sensors="snapshot.temperature" />
        </UCard>

        <UCard>
          <template #header>
            <h2 class="font-medium">
              Disks
            </h2>
          </template>
          <DiskTable :disks="snapshot.disks" />
        </UCard>
      </div>

      <HostCard :host="snapshot.host" />
    </template>
  </div>
</template>
