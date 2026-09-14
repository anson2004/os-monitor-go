package collector

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/mem"
)

// MemoryInfo covers physical RAM and swap.
type MemoryInfo struct {
	TotalBytes      uint64  `json:"total_bytes"`
	UsedBytes       uint64  `json:"used_bytes"`
	AvailableBytes  uint64  `json:"available_bytes"`
	UsedPercent     float64 `json:"used_percent"`
	SwapTotalBytes  uint64  `json:"swap_total_bytes"`
	SwapUsedBytes   uint64  `json:"swap_used_bytes"`
	SwapUsedPercent float64 `json:"swap_used_percent"`
}

func (c *Collector) memory(ctx context.Context, snap *Snapshot) error {
	vm, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return fmt.Errorf("memory: %w", err)
	}

	snap.Memory = MemoryInfo{
		TotalBytes:     vm.Total,
		UsedBytes:      vm.Used,
		AvailableBytes: vm.Available,
		UsedPercent:    round2(vm.UsedPercent),
	}

	if sw, err := mem.SwapMemoryWithContext(ctx); err == nil {
		snap.Memory.SwapTotalBytes = sw.Total
		snap.Memory.SwapUsedBytes = sw.Used
		snap.Memory.SwapUsedPercent = round2(sw.UsedPercent)
	}
	return nil
}
