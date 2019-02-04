// spot.go
package gospot

import "C"

import (
	"errors"
	"fmt"
    "math"
)

var (
	spot_new               func(float64, int32, float64, uint8, uint8, uint8, uint8, int32) uintptr
	spot_delete            func(uintptr)
	spot_step              func(uintptr, float64) int32
	spot_getUpperThreshold func(uintptr) float64
	spot_getLowerThreshold func(uintptr) float64
	spot_getUpper_t        func(uintptr) float64
	spot_getLower_t        func(uintptr) float64
	spot_set_q             func(uintptr, float64)
	spot_up_probability    func(uintptr, float64) float64
	spot_down_probability  func(uintptr, float64) float64
)

// LoadSymbolsSpot It loads the symbols related to the Spot object from
// the C++ library libspot. It returns an error if a loading fails, nil
// pointer otherwise
func LoadSymbolsSpot() error {
	var err error
	err = libspot.Sym("Spot_new", &spot_new)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading Spot_new (%s)", err.Error()))
	}
	err = libspot.Sym("Spot_delete", &spot_delete)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading Spot_delete (%s)", err.Error()))
	}
	err = libspot.Sym("Spot_step", &spot_step)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading Spot_step (%s)", err.Error()))
	}
	err = libspot.Sym("Spot_getUpperThreshold", &spot_getUpperThreshold)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading Spot_getUpperThreshold (%s)", err.Error()))
	}
	err = libspot.Sym("Spot_getLowerThreshold", &spot_getLowerThreshold)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading Spot_getLowerThreshold (%s)", err.Error()))
	}
	err = libspot.Sym("Spot_getUpper_t", &spot_getUpper_t)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading Spot_getUpper_t (%s)", err.Error()))
	}
	err = libspot.Sym("Spot_getLower_t", &spot_getLower_t)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading Spot_getLower_t (%s)", err.Error()))
	}
	err = libspot.Sym("Spot_set_q", &spot_set_q)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading Spot_set_q (%s)", err.Error()))
	}
	err = libspot.Sym("Spot_up_probability", &spot_up_probability)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading Spot_up_probability (%s)", err.Error()))
	}
	err = libspot.Sym("Spot_down_probability", &spot_down_probability)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading Spot_down_probability (%s)", err.Error()))
	}
	return nil
}

// Spot This object embeds a pointer to a C++ object Spot
type Spot struct {
	// ptr pointer to Spot instance
	ptr uintptr
	// up/down flag
	Up bool
	Down bool
}

func NewDefaultSpot() *Spot {
	return NewSpotFromConfig(DefaultSpotConfig)
}

// NewSpot Constructor
func NewSpot(q float64, n_init int32, level float64, up bool, down bool, alert bool, bounded bool, max_excess int32) *Spot {
	up_u8 := bool2uint8(up)
	down_u8 := bool2uint8(down)
	alert_u8 := bool2uint8(alert)
	bounded_u8 := bool2uint8(bounded)
	return &Spot{ptr: spot_new(q, n_init, level, up_u8, down_u8, alert_u8, bounded_u8, max_excess), Up: up, Down: down}
}

func NewSpotFromConfig(sc SpotConfig) *Spot {
	return NewSpot(sc.Q, sc.N_init, sc.Level, sc.Up, sc.Down, sc.Alert, sc.Bounded, sc.Max_excess)
}

// Delete Destructor
func (s *Spot) Delete() {
	spot_delete(s.ptr)
}

// Step Method which update the Spot instance according to a new incoming value
func (s *Spot) Step(x float64) int32 {
	return spot_step(s.ptr, x)
}

// GetUpperT Returns the upper threshold t
func (s *Spot) GetUpperT() float64 {
    if s.Up {
		return spot_getUpper_t(s.ptr)
    } else {
        return math.NaN()
    }
}

// GetUpperT Returns the lower threshold t
func (s *Spot) GetLowerT() float64 {
    if s.Down {
		return spot_getLower_t(s.ptr)
    } else {
        return math.NaN()
    }
}

// GetUpperT Returns the upper decision threshold
func (s *Spot) GetUpperThreshold() float64 {
    if s.Up {
		return spot_getUpperThreshold(s.ptr)
    } else {
        return math.NaN()
    }
}

// GetUpperT Returns the lower decision threshold
func (s *Spot) GetLowerThreshold() float64 {
    if s.Down {
		return spot_getLowerThreshold(s.ptr)
    } else {
        return math.NaN()
    }
}

// SetQ Change the value of the decision probability.
// It then changes the decision thresholds
func (s *Spot) SetQ(q float64) {
	spot_set_q(s.ptr, q)
}

// UpProbability Given a quantile z, computes the probability
// to observe a value higher than z
func (s *Spot) UpProbability(z float64) float64 {
    if s.Up {
		return spot_up_probability(s.ptr, z)
    } else {
        return math.NaN()
    }
}

// DownProbability Given a quantile z, computes the probability
// to observe a value lower than z
func (s *Spot) DownProbability(z float64) float64 {
    if s.Down {
		return spot_down_probability(s.ptr, z)
    } else {
        return math.NaN()
    }
}

func (s *Spot) Status() SpotStatus {
	status_ptr := spot_status_new(s.ptr)
	defer spot_status_delete(status_ptr)
	return SpotStatus{
		N:       status_get_n(status_ptr),
		Ex_up:   status_get_ex_up(status_ptr),
		Ex_down: status_get_ex_down(status_ptr),
		Nt_up:   status_get_Nt_up(status_ptr),
		Nt_down: status_get_Nt_down(status_ptr),
		Al_up:   status_get_al_up(status_ptr),
		Al_down: status_get_al_down(status_ptr),
		T_up:    status_get_t_up(status_ptr),
		T_down:  status_get_t_down(status_ptr),
		Z_up:    status_get_z_up(status_ptr),
		Z_down:  status_get_z_down(status_ptr)}
}

func (s *Spot) Config() SpotConfig {
	config_ptr := spot_config_new(s.ptr)
	defer spot_config_delete(config_ptr)
	return SpotConfig{
		Q:          config_get_q(config_ptr),
		N_init:     config_get_n_init(config_ptr),
		Level:      config_get_level(config_ptr),
		Up:         config_get_up(config_ptr) == 1,
		Down:       config_get_down(config_ptr) == 1,
		Alert:      config_get_alert(config_ptr) == 1,
		Bounded:    config_get_bounded(config_ptr) == 1,
		Max_excess: config_get_max_excess(config_ptr)}
}
