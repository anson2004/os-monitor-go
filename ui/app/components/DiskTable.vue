<script setup lang="ts">
import type { TableColumn } from '@nuxt/ui'
import type { DiskInfo } from '#shared/types/metrics'
import { formatBytes, usageColor } from '~/utils/format'

const props = defineProps<{ disks: DiskInfo[] }>()

// Hide tiny pseudo filesystems and de-duplicate volumes that share one pool
// (APFS volumes in the same container report identical totals).
const rows = computed(() => {
  const seen = new Set<string>()
  return props.disks
    .filter(d => d.total_bytes > 1024 * 1024 * 1024)
    .filter((d) => {
      const key = `${d.total_bytes}:${d.used_bytes}:${d.free_bytes}`
      if (seen.has(key)) return false
      seen.add(key)
      return true
    })
    .sort((a, b) => b.used_percent - a.used_percent)
})

const UProgress = resolveComponent('UProgress')

const columns: TableColumn<DiskInfo>[] = [
  { accessorKey: 'mount_point', header: 'Mount' },
  { accessorKey: 'fs_type', header: 'Type' },
  {
    accessorKey: 'used_bytes',
    header: 'Used',
    cell: ({ row }) => `${formatBytes(row.original.used_bytes)} / ${formatBytes(row.original.total_bytes)}`
  },
  {
    accessorKey: 'used_percent',
    header: 'Usage',
    cell: ({ row }) =>
      h('div', { class: 'flex items-center gap-2 min-w-32' }, [
        h(UProgress, {
          modelValue: row.original.used_percent,
          color: usageColor(row.original.used_percent),
          size: 'sm',
          class: 'flex-1'
        }),
        h('span', { class: 'w-12 text-right tabular-nums' }, `${row.original.used_percent.toFixed(0)}%`)
      ])
  }
]
</script>

<template>
  <UEmpty
    v-if="!rows.length"
    icon="i-lucide-hard-drive"
    title="No disks reported"
  />
  <UTable
    v-else
    :data="rows"
    :columns="columns"
    class="max-h-80"
  />
</template>
