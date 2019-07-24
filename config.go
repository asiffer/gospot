// status.go

package gospot

import (
	"fmt"
)

// SpotConfig is the structure embedding the Spot configuration
type SpotConfig struct {
	// the main parameter ( P(X>zQ) < q )
	Q float64 `json:"q"`
	// number of observation to perform calibration
	Ninit int `json:"n_init"`
	// level of the update threshold (0<l<1)
	Level float64 `json:"level"`
	// if true, compute upper threshold
	Up bool `json:"up"`
	// if true, compute lower threshold
	Down bool `json:"down"`
	// if true, the algorithm triggers alarms (the outlier is not taking into account in the model)
	Alert bool `json:"alert"`
	// if true, the number of stored will be bounded by max_excess
	Bounded bool `json:"bounded"`
	// Maximum number of stored excesses (bounded mode)
	MaxExcess int `json:"max_excess"`
}

// DSpotConfig is the structure embedding the DSpot config (SpotConfig + depth)
type DSpotConfig struct {
	SpotConfig
	// Depth is the size of the underlying moving average
	Depth int
}

func (sc SpotConfig) String() string {
	return fmt.Sprintf("%10s %.6f\n%10s %d\n%10s %.6f\n%10s %t\n%10s %t\n%10s %t\n%10s %t\n%10s %d\n",
		"q", sc.Q, "n_init", sc.Ninit, "level", sc.Level, "up", sc.Up, "down", sc.Down, "alert", sc.Alert, "bounded", sc.Bounded, "max_excess", sc.MaxExcess)
}
