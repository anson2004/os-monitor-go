<script setup lang="ts">
const props = defineProps<{ perCore: number[] }>()

const data = computed(() =>
  props.perCore.map((usage, i) => ({ core: `${i}`, usage }))
)

const categories: Record<string, BulletLegendItemInterface> = {
  usage: { name: 'Usage %', color: '#3b82f6' }
}

const xFormatter = (i: number) => `C${data.value[i]?.core ?? ''}`
const yFormatter = (v: number) => `${v}%`
</script>

<template>
  <BarChart
    :data="data"
    :categories="categories"
    :y-axis="['usage']"
    :height="260"
    :y-domain="[0, 100]"
    :y-num-ticks="5"
    :x-num-ticks="data.length"
    :x-formatter="xFormatter"
    :y-formatter="yFormatter"
    :radius="3"
    :hide-legend="true"
    :y-grid-line="true"
  />
</template>
