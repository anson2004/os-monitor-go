// Package collector gathers operating-system metrics into a single Snapshot.
package collector

import (
	"context"
	"math"
	"time"
)

// Snapshot is one point-in-time reading of the machine.
type Snapshot struct {
	Timestamp   time.Time     `json:"timestamp"`
	Host        HostInfo      `json:"host"`
	CPU         CPUInfo       `json:"cpu"`
	Memory      MemoryInfo    `json:"memory"`
	Disks       []DiskInfo    `json:"disks"`
	Temperature []Temperature `json:"temperature"`
	Errors      []string      `json:"errors,omitempty"`
}

// Collector produces Snapshots.
type Collector struct {
	cpuSample time.Duration
}

// New returns a Collector with default settings.
func New() *Collector {
	return &Collector{cpuSample: 500 * time.Millisecond}
}

// step fills one part of a Snapshot. Steps are independent: a failing source
// (for example, no temperature sensors inside a container) is reported in
// Snapshot.Errors without affecting the other groups.
type step func(ctx context.Context, snap *Snapshot) error

// Collect runs every step and returns the combined Snapshot.
func (c *Collector) Collect(ctx context.Context) Snapshot {
	snap := Snapshot{Timestamp: time.Now().UTC()}

	steps := []step{
		c.host,
		c.cpu,
		c.memory,
		c.disks,
		c.temperature,
	}
	for _, run := range steps {
		if err := run(ctx, &snap); err != nil {
			snap.Errors = append(snap.Errors, err.Error())
		}
	}
	return snap
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
