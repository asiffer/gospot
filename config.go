// status.go

package gospot

import (
	"fmt"
)

// SpotConfig is the structure embedding the Spot configuration
type SpotConfig struct {
	// the main parameter ( P(X>zQ) < q )
	Q float64
	// number of observation to perform calibration
	Ninit int
	// level of the update threshold (0<l<1)
	Level float64
	// if true, compute upper threshold
	Up bool
	// if true, compute lower threshold
	Down bool
	// if true, the algorithm triggers alarms (the outlier is not taking into account in the model)
	Alert bool
	// if true, the number of stored will be bounded by max_excess
	Bounded bool
	// Maximum number of stored excesses (bounded mode)
	MaxExcess int
}

// DSpotConfig is the structure embedding the DSpot config (SpotConfig + depth)
type DSpotConfig struct {
	SpotConfig
	// Depth is the size of the underlying moving average
	Depth int
}

// var (
// 	spotConfigNew      func(uintptr) uintptr
// 	spotConfigDelete   func(uintptr)
// 	configGetQ         func(uintptr) float64
// 	configGetBounded   func(uintptr) int32
// 	configGetMaxExcess func(uintptr) int32
// 	configGetAlert     func(uintptr) int32
// 	configGetUp        func(uintptr) int32
// 	configGetDown      func(uintptr) int32
// 	configGetNinit     func(uintptr) int32
// 	configGetLevel     func(uintptr) float64
// 	// DefaultSpotConfig is the default structure to create a Spot Object
// 	DefaultSpotConfig SpotConfig
// 	// DefaultDSpotConfig is the default structure to create a DSpot Object
// 	DefaultDSpotConfig DSpotConfig
// )

// var configSymbols = map[string]interface{}{
// 	"Spot_config_ptr":        &spotConfigNew,
// 	"Spot_config_delete":     &spotConfigDelete,
// 	"_config_get_q":          &configGetQ,
// 	"_config_get_bounded":    &configGetBounded,
// 	"_config_get_max_excess": &configGetMaxExcess,
// 	"_config_get_alert":      &configGetAlert,
// 	"_config_get_up":         &configGetUp,
// 	"_config_get_down":       &configGetDown,
// 	"_config_get_n_init":     &configGetNinit,
// 	"_config_get_level":      &configGetLevel,
// }

// // LoadSymbolsConfig loads the symbols related to the SpotStatus object
// // from the C++ library libspot. It returns an error if a loading fails, nil
// // pointer otherwise
// func LoadSymbolsConfig() error {
// 	var err error
// 	for k, f := range configSymbols {
// 		err = libspot.Sym(k, f)
// 		if err != nil {
// 			return fmt.Errorf("Error in loading %s (%s)", k, err.Error())
// 		}
// 	}
// 	return nil
// }

func (sc SpotConfig) String() string {
	return fmt.Sprintf("%10s %.6f\n%10s %d\n%10s %.6f\n%10s %t\n%10s %t\n%10s %t\n%10s %t\n%10s %d\n",
		"q", sc.Q, "n_init", sc.Ninit, "level", sc.Level, "up", sc.Up, "down", sc.Down, "alert", sc.Alert, "bounded", sc.Bounded, "max_excess", sc.MaxExcess)
}
