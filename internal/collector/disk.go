package collector

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/disk"
)

// DiskInfo is the usage of one mounted filesystem.
type DiskInfo struct {
	Device      string  `json:"device"`
	MountPoint  string  `json:"mount_point"`
	FSType      string  `json:"fs_type"`
	TotalBytes  uint64  `json:"total_bytes"`
	UsedBytes   uint64  `json:"used_bytes"`
	FreeBytes   uint64  `json:"free_bytes"`
	UsedPercent float64 `json:"used_percent"`
}

func (c *Collector) disks(ctx context.Context, snap *Snapshot) error {
	parts, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		return fmt.Errorf("disks: %w", err)
	}

	snap.Disks = make([]DiskInfo, 0, len(parts))
	for _, p := range parts {
		usage, err := disk.UsageWithContext(ctx, p.Mountpoint)
		if err != nil || usage.Total == 0 {
			continue // pseudo filesystems and unreadable mounts
		}
		snap.Disks = append(snap.Disks, DiskInfo{
			Device:      p.Device,
			MountPoint:  p.Mountpoint,
			FSType:      p.Fstype,
			TotalBytes:  usage.Total,
			UsedBytes:   usage.Used,
			FreeBytes:   usage.Free,
			UsedPercent: round2(usage.UsedPercent),
		})
	}
	return nil
}
