import type { HistoryPoint, Snapshot } from '#shared/types/metrics'

const MAX_HISTORY = 60

/**
 * Polls the Go API for the latest snapshot and keeps a rolling history of
 * CPU and memory usage for trend charts. Client-side only.
 */
export function useMetrics(intervalMs = 3000) {
  const snapshot = ref<Snapshot | null>(null)
  const history = ref<HistoryPoint[]>([])
  const error = ref<string | null>(null)
  const loading = ref(true)
  const lastUpdated = ref<Date | null>(null)

  let timer: ReturnType<typeof setInterval> | undefined

  async function refresh() {
    try {
      const snap = await $fetch<Snapshot>('/api/metrics')
      // The server caches its snapshot, so consecutive polls may return the
      // same reading. Only append to history when the timestamp changes.
      if (snap.timestamp !== snapshot.value?.timestamp) {
        history.value = [
          ...history.value,
          {
            time: snap.timestamp,
            cpu: snap.cpu.usage_percent,
            memory: snap.memory.used_percent,
            swap: snap.memory.swap_used_percent
          }
        ].slice(-MAX_HISTORY)
      }
      snapshot.value = snap
      lastUpdated.value = new Date()
      error.value = null
    } catch (e) {
      error.value = e instanceof Error ? e.message : String(e)
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    refresh()
    timer = setInterval(refresh, intervalMs)
  })

  onUnmounted(() => {
    if (timer) clearInterval(timer)
  })

  return { snapshot, history, error, loading, lastUpdated, refresh }
}
