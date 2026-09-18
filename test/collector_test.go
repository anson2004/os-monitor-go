package test

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"osmonitor/internal/collector"
)

// These tests run the real collector against the machine executing them, so
// they assert on properties every host must satisfy rather than exact values.

func TestCollectReturnsPopulatedSnapshot(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	snap := collector.New().Collect(ctx)

	if snap.Timestamp.IsZero() {
		t.Fatal("timestamp not set")
	}
	if since := time.Since(snap.Timestamp); since < 0 || since > time.Minute {
		t.Errorf("timestamp %v is not recent (now-ts = %v)", snap.Timestamp, since)
	}

	if snap.CPU.LogicalCores <= 0 {
		t.Errorf("logical cores = %d, want > 0", snap.CPU.LogicalCores)
	}
	if got, want := len(snap.CPU.PerCore), snap.CPU.LogicalCores; got != want {
		t.Errorf("per-core samples = %d, want %d (one per logical core)", got, want)
	}
	assertPercent(t, "cpu usage", snap.CPU.UsagePercent)
	for i, v := range snap.CPU.PerCore {
		assertPercent(t, fmt.Sprintf("cpu core %d", i), v)
	}

	if snap.Memory.TotalBytes == 0 {
		t.Error("memory total is zero")
	}
	if snap.Memory.UsedBytes > snap.Memory.TotalBytes {
		t.Errorf("memory used %d exceeds total %d", snap.Memory.UsedBytes, snap.Memory.TotalBytes)
	}
	assertPercent(t, "memory usage", snap.Memory.UsedPercent)

	if snap.Host.Hostname == "" {
		t.Error("hostname is empty")
	}
	if snap.Host.OS == "" {
		t.Error("host OS is empty")
	}

	for _, d := range snap.Disks {
		if d.TotalBytes == 0 {
			t.Errorf("disk %s reported with zero total; zero-size mounts should be skipped", d.MountPoint)
		}
		if d.MountPoint == "" {
			t.Errorf("disk %s has empty mount point", d.Device)
		}
		assertPercent(t, "disk "+d.MountPoint, d.UsedPercent)
	}

	// Sensors may legitimately be absent (containers, VMs), but any reading
	// that is present must have survived the plausibility filter.
	for _, s := range snap.Temperature {
		if s.Sensor == "" {
			t.Error("temperature reading with empty sensor name")
		}
		if s.Celsius <= -50 || s.Celsius >= 200 {
			t.Errorf("sensor %s = %.2f°C, outside plausible range", s.Sensor, s.Celsius)
		}
	}
}

func TestCollectRoundsToTwoDecimals(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	snap := collector.New().Collect(ctx)

	values := map[string]float64{
		"cpu usage":    snap.CPU.UsagePercent,
		"memory usage": snap.Memory.UsedPercent,
		"swap usage":   snap.Memory.SwapUsedPercent,
		"load 1":       snap.Host.Load1,
	}
	for name, v := range values {
		if diff := math.Abs(v*100 - math.Round(v*100)); diff > 1e-6 {
			t.Errorf("%s = %v is not rounded to two decimals", name, v)
		}
	}
}

func TestCollectWithCancelledContextStillReturns(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan collector.Snapshot, 1)
	go func() { done <- collector.New().Collect(ctx) }()

	select {
	case snap := <-done:
		if snap.Timestamp.IsZero() {
			t.Error("timestamp should be set even when collection is cancelled")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Collect did not return promptly with a cancelled context")
	}
}

func assertPercent(t *testing.T, name string, v float64) {
	t.Helper()
	if v < 0 || v > 100 {
		t.Errorf("%s = %v, want 0..100", name, v)
	}
}
