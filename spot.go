// spot.go

package gospot

import "C"

import (
	"fmt"
	"math"
)

var (
	spotNew               func(float64, int32, float64, uint8, uint8, uint8, uint8, int32) uintptr
	spotDelete            func(uintptr)
	spotStep              func(uintptr, float64) int32
	spotGetUpperThreshold func(uintptr) float64
	spotGetLowerThreshold func(uintptr) float64
	spotGetUpperT         func(uintptr) float64
	spotGetLowerT         func(uintptr) float64
	spotSetQ              func(uintptr, float64)
	spotUpProbability     func(uintptr, float64) float64
	spotDownProbability   func(uintptr, float64) float64
)

// LoadSymbolsSpot It loads the symbols related to the Spot object from
// the C++ library libspot. It returns an error if a loading fails, nil
// pointer otherwise
func LoadSymbolsSpot() error {
	var err error
	err = libspot.Sym("Spot_new", &spotNew)
	if err != nil {
		return fmt.Errorf("Error in loading Spot_new (%s)", err.Error())
	}
	err = libspot.Sym("Spot_delete", &spotDelete)
	if err != nil {
		return fmt.Errorf("Error in loading Spot_delete (%s)", err.Error())
	}
	err = libspot.Sym("Spot_step", &spotStep)
	if err != nil {
		return fmt.Errorf("Error in loading Spot_step (%s)", err.Error())
	}
	err = libspot.Sym("Spot_getUpperThreshold", &spotGetUpperThreshold)
	if err != nil {
		return fmt.Errorf("Error in loading Spot_getUpperThreshold (%s)", err.Error())
	}
	err = libspot.Sym("Spot_getLowerThreshold", &spotGetLowerThreshold)
	if err != nil {
		return fmt.Errorf("Error in loading Spot_getLowerThreshold (%s)", err.Error())
	}
	err = libspot.Sym("Spot_getUpper_t", &spotGetUpperT)
	if err != nil {
		return fmt.Errorf("Error in loading Spot_getUpper_t (%s)", err.Error())
	}
	err = libspot.Sym("Spot_getLower_t", &spotGetLowerT)
	if err != nil {
		return fmt.Errorf("Error in loading Spot_getLower_t (%s)", err.Error())
	}
	err = libspot.Sym("Spot_set_q", &spotSetQ)
	if err != nil {
		return fmt.Errorf("Error in loading Spot_set_q (%s)", err.Error())
	}
	err = libspot.Sym("Spot_up_probability", &spotUpProbability)
	if err != nil {
		return fmt.Errorf("Error in loading Spot_up_probability (%s)", err.Error())
	}
	err = libspot.Sym("Spot_down_probability", &spotDownProbability)
	if err != nil {
		return fmt.Errorf("Error in loading Spot_down_probability (%s)", err.Error())
	}
	return nil
}

// Spot This object embeds a pointer to a C++ object Spot
type Spot struct {
	// ptr pointer to Spot instance
	ptr uintptr
	// up/down flag
	Up   bool
	Down bool
}

// NewDefaultSpot is the default constructor of a Spot object
func NewDefaultSpot() *Spot {
	return NewSpotFromConfig(DefaultSpotConfig)
}

// NewSpot is the Spot constructor
func NewSpot(q float64, nInit int32, level float64, up bool, down bool, alert bool, bounded bool, maxExcess int32) *Spot {
	up8 := bool2uint8(up)
	down8 := bool2uint8(down)
	alert8 := bool2uint8(alert)
	bounded8 := bool2uint8(bounded)
	return &Spot{ptr: spotNew(q, nInit, level, up8, down8, alert8, bounded8, maxExcess), Up: up, Down: down}
}

// NewSpotFromConfig creates a Spot instance from a config structure
func NewSpotFromConfig(sc SpotConfig) *Spot {
	return NewSpot(sc.Q, sc.Ninit, sc.Level, sc.Up, sc.Down, sc.Alert, sc.Bounded, sc.MaxExcess)
}

// Delete Destructor
func (s *Spot) Delete() {
	spotDelete(s.ptr)
}

// Step Method which update the Spot instance according to a new incoming value
func (s *Spot) Step(x float64) int32 {
	return spotStep(s.ptr, x)
}

// GetUpperT Returns the upper threshold t
func (s *Spot) GetUpperT() float64 {
	if s.Up {
		return spotGetUpperT(s.ptr)
	}
	return math.NaN()

}

// GetLowerT Returns the lower threshold t
func (s *Spot) GetLowerT() float64 {
	if s.Down {
		return spotGetLowerT(s.ptr)
	}
	return math.NaN()

}

// GetUpperThreshold returns the upper decision threshold
func (s *Spot) GetUpperThreshold() float64 {
	if s.Up {
		return spotGetUpperThreshold(s.ptr)
	}
	return math.NaN()

}

// GetLowerThreshold returns the lower decision threshold
func (s *Spot) GetLowerThreshold() float64 {
	if s.Down {
		return spotGetLowerThreshold(s.ptr)
	}
	return math.NaN()

}

// SetQ Change the value of the decision probability.
// It then changes the decision thresholds
func (s *Spot) SetQ(q float64) {
	spotSetQ(s.ptr, q)
}

// UpProbability Given a quantile z, computes the probability
// to observe a value greater than z
func (s *Spot) UpProbability(z float64) float64 {
	if s.Up {
		return spotUpProbability(s.ptr, z)
	}
	return math.NaN()

}

// DownProbability Given a quantile z, computes the probability
// to observe a value lower than z
func (s *Spot) DownProbability(z float64) float64 {
	if s.Down {
		return spotDownProbability(s.ptr, z)
	}
	return math.NaN()

}

// Status returns the current status of the Spot instance
func (s *Spot) Status() SpotStatus {
	statusPtr := spotStatusNew(s.ptr)
	defer spotStatusDelete(statusPtr)
	return SpotStatus{
		N:      statusGetN(statusPtr),
		ExUp:   statusGetExUp(statusPtr),
		ExDown: statusGetExDown(statusPtr),
		NtUp:   statusGetNtUp(statusPtr),
		NtDown: statusGetNtDown(statusPtr),
		AlUp:   statusGetAlUp(statusPtr),
		AlDown: statusGetAlDown(statusPtr),
		TUp:    statusGetTUp(statusPtr),
		TDown:  statusGetTDown(statusPtr),
		ZUp:    statusGetZUp(statusPtr),
		ZDown:  statusGetZDown(statusPtr)}
}

// Config returns the configuration of the Spot instance
func (s *Spot) Config() SpotConfig {
	configPtr := spotConfigNew(s.ptr)
	defer spotConfigDelete(configPtr)
	return SpotConfig{
		Q:         configGetQ(configPtr),
		Ninit:     configGetNinit(configPtr),
		Level:     configGetLevel(configPtr),
		Up:        configGetUp(configPtr) == 1,
		Down:      configGetDown(configPtr) == 1,
		Alert:     configGetAlert(configPtr) == 1,
		Bounded:   configGetBounded(configPtr) == 1,
		MaxExcess: configGetMaxExcess(configPtr)}
}
