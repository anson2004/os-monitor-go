<script setup lang="ts">
import { usageColor } from '~/utils/format'

defineProps<{
  title: string
  icon: string
  value: string
  subtitle?: string
  /** When given, a progress bar coloured by utilisation is shown. */
  percent?: number
}>()
</script>

<template>
  <UCard>
    <div class="flex items-start justify-between gap-2">
      <div class="min-w-0">
        <p class="text-sm text-muted">
          {{ title }}
        </p>
        <p class="mt-1 text-3xl font-semibold tabular-nums">
          {{ value }}
        </p>
      </div>
      <UIcon
        :name="icon"
        class="size-6 shrink-0 text-muted"
      />
    </div>

    <UProgress
      v-if="percent !== undefined"
      :model-value="percent"
      :color="usageColor(percent)"
      size="sm"
      class="mt-3"
    />

    <p
      v-if="subtitle"
      class="mt-2 truncate text-xs text-muted"
      :title="subtitle"
    >
      {{ subtitle }}
    </p>
  </UCard>
</template>
