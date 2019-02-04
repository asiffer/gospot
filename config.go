// status.go
package gospot

import (
	"errors"
	"fmt"
)

type SpotConfig struct {
	// the main parameter ( P(X>z_q) < q )
	Q float64
	// number of observation to perform calibration
	N_init int32
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
	Max_excess int32
}

type DSpotConfig struct {
	SpotConfig
	// Depth is the size of the underlying moving average
	Depth int
}

var (
	spot_config_new       func(uintptr) uintptr
	spot_config_delete    func(uintptr)
	config_get_q          func(uintptr) float64
	config_get_bounded    func(uintptr) int32
	config_get_max_excess func(uintptr) int32
	config_get_alert      func(uintptr) int32
	config_get_up         func(uintptr) int32
	config_get_down       func(uintptr) int32
	config_get_n_init     func(uintptr) int32
	config_get_level      func(uintptr) float64
	DefaultSpotConfig     SpotConfig
	DefaultDSpotConfig    DSpotConfig
)

// LoadSymbolsStatus It loads the symbols related to the SpotStatus object
// from the C++ library libspot. It returns an error if a loading fails, nil
// pointer otherwise
func LoadSymbolsConfig() error {
	var err error
	err = libspot.Sym("Spot_config_ptr", &spot_config_new)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading Spot_config_ptr(%s)", err.Error()))
	}
	err = libspot.Sym("Spot_config_delete", &spot_config_delete)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading Spot_config_delete (%s)", err.Error()))
	}
	err = libspot.Sym("_config_get_q", &config_get_q)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading config_get_q (%s)", err.Error()))
	}
	err = libspot.Sym("_config_get_bounded", &config_get_bounded)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading config_get_bounded (%s)", err.Error()))
	}
	err = libspot.Sym("_config_get_max_excess", &config_get_max_excess)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading config_get_max_excess (%s)", err.Error()))
	}
	err = libspot.Sym("_config_get_alert", &config_get_alert)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading config_get_alert (%s)", err.Error()))
	}
	err = libspot.Sym("_config_get_up", &config_get_up)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading config_get_up (%s)", err.Error()))
	}
	err = libspot.Sym("_config_get_down", &config_get_down)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading config_get_down (%s)", err.Error()))
	}
	err = libspot.Sym("_config_get_n_init", &config_get_n_init)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading config_get_n_init (%s)", err.Error()))
	}
	err = libspot.Sym("_config_get_level", &config_get_level)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading config_get_level (%s)", err.Error()))
	}
	return nil
}

func (sc SpotConfig) String() string {
	return fmt.Sprintf("%10s %.6f\n%10s %d\n%10s %.6f\n%10s %t\n%10s %t\n%10s %t\n%10s %t\n%10s %d\n",
		"q", sc.Q, "n_init", sc.N_init, "level", sc.Level, "up", sc.Up, "down", sc.Down, "alert", sc.Alert, "bounded", sc.Bounded, "max_excess", sc.Max_excess)
}
