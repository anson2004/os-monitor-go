// Mirrors the JSON produced by the Go collector (internal/collector).

export interface HostInfo {
  hostname: string
  os: string
  platform: string
  platform_version: string
  kernel_version: string
  arch: string
  uptime_seconds: number
  processes: number
  load_1: number
  load_5: number
  load_15: number
}

export interface CPUInfo {
  model_name: string
  physical_cores: number
  logical_cores: number
  frequency_mhz: number
  usage_percent: number
  per_core_percent: number[]
}

export interface MemoryInfo {
  total_bytes: number
  used_bytes: number
  available_bytes: number
  used_percent: number
  swap_total_bytes: number
  swap_used_bytes: number
  swap_used_percent: number
}

export interface DiskInfo {
  device: string
  mount_point: string
  fs_type: string
  total_bytes: number
  used_bytes: number
  free_bytes: number
  used_percent: number
}

export interface Temperature {
  sensor: string
  celsius: number
  high?: number
  critical?: number
}

export interface Snapshot {
  timestamp: string
  host: HostInfo
  cpu: CPUInfo
  memory: MemoryInfo
  disks: DiskInfo[]
  temperature: Temperature[]
  errors?: string[]
}

/** One point in the client-side history used for trend charts. */
export interface HistoryPoint {
  time: string
  cpu: number
  memory: number
  swap: number
}
