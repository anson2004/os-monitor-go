package collector

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/cpu"
)

// CPUInfo describes processor layout and current utilisation.
type CPUInfo struct {
	ModelName     string    `json:"model_name"`
	PhysicalCores int       `json:"physical_cores"`
	LogicalCores  int       `json:"logical_cores"`
	FrequencyMHz  float64   `json:"frequency_mhz"`
	UsagePercent  float64   `json:"usage_percent"`
	PerCore       []float64 `json:"per_core_percent"`
}

func (c *Collector) cpu(ctx context.Context, snap *Snapshot) error {
	info := &snap.CPU

	if infos, err := cpu.InfoWithContext(ctx); err == nil && len(infos) > 0 {
		info.ModelName = infos[0].ModelName
		info.FrequencyMHz = infos[0].Mhz
	}
	if n, err := cpu.CountsWithContext(ctx, false); err == nil {
		info.PhysicalCores = n
	}
	if n, err := cpu.CountsWithContext(ctx, true); err == nil {
		info.LogicalCores = n
	}

	// Sample per-core once and derive the total, instead of blocking for the
	// sample interval twice.
	perCore, err := cpu.PercentWithContext(ctx, c.cpuSample, true)
	if err != nil {
		return fmt.Errorf("cpu usage: %w", err)
	}
	if len(perCore) == 0 {
		return nil
	}

	var sum float64
	info.PerCore = make([]float64, len(perCore))
	for i, v := range perCore {
		info.PerCore[i] = round2(v)
		sum += v
	}
	info.UsagePercent = round2(sum / float64(len(perCore)))
	return nil
}
