package collector

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/sensors"
)

// Temperature is one sensor reading in degrees Celsius.
type Temperature struct {
	Sensor   string  `json:"sensor"`
	Celsius  float64 `json:"celsius"`
	High     float64 `json:"high,omitempty"`
	Critical float64 `json:"critical,omitempty"`
}

// temperature reads hardware sensors. Availability depends on the platform:
// Linux uses /sys/class/hwmon and /sys/class/thermal, macOS reads the SMC,
// and containers see nothing unless the host /sys is mounted.
func (c *Collector) temperature(ctx context.Context, snap *Snapshot) error {
	readings, err := sensors.TemperaturesWithContext(ctx)
	if err != nil && len(readings) == 0 {
		return fmt.Errorf("temperature: %w", err)
	}

	snap.Temperature = make([]Temperature, 0, len(readings))
	for _, r := range readings {
		if !plausible(r.Temperature) {
			continue
		}
		snap.Temperature = append(snap.Temperature, Temperature{
			Sensor:   r.SensorKey,
			Celsius:  round2(r.Temperature),
			High:     round2(r.High),
			Critical: round2(r.Critical),
		})
	}
	return nil
}

// plausible rejects sentinel values (for example -9200) that some platforms
// report for unpopulated sensor slots.
func plausible(celsius float64) bool {
	return celsius > -50 && celsius < 200
}
