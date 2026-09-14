package collector

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
)

// HostInfo identifies the machine and reports uptime and load.
type HostInfo struct {
	Hostname        string  `json:"hostname"`
	OS              string  `json:"os"`
	Platform        string  `json:"platform"`
	PlatformVersion string  `json:"platform_version"`
	KernelVersion   string  `json:"kernel_version"`
	Arch            string  `json:"arch"`
	UptimeSeconds   uint64  `json:"uptime_seconds"`
	Processes       uint64  `json:"processes"`
	Load1           float64 `json:"load_1"`
	Load5           float64 `json:"load_5"`
	Load15          float64 `json:"load_15"`
}

func (c *Collector) host(ctx context.Context, snap *Snapshot) error {
	h, err := host.InfoWithContext(ctx)
	if err != nil {
		return fmt.Errorf("host: %w", err)
	}

	snap.Host = HostInfo{
		Hostname:        h.Hostname,
		OS:              h.OS,
		Platform:        h.Platform,
		PlatformVersion: h.PlatformVersion,
		KernelVersion:   h.KernelVersion,
		Arch:            h.KernelArch,
		UptimeSeconds:   h.Uptime,
		Processes:       h.Procs,
	}

	if avg, err := load.AvgWithContext(ctx); err == nil {
		snap.Host.Load1 = round2(avg.Load1)
		snap.Host.Load5 = round2(avg.Load5)
		snap.Host.Load15 = round2(avg.Load15)
	}
	return nil
}
