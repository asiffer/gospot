// spot.go

package gospot

import (
	"math"
	"sort"
)

const (
	// Normal refers to normal data
	Normal = 0
	// AlertUp refers to an upper anomaly
	AlertUp = 1
	// AlertDown refers to lower anomaly
	AlertDown = -1
	// ExcessUp refers to a data used to the upper tail fit
	ExcessUp = 2
	// ExcessDown refers to a data used to the lower tail fit
	ExcessDown = -2
	// InitBatch refers to a data stored in the initial batch (before calibration)
	InitBatch = 3
	// Calibration refers to the last InitBatch data (the calibration step is performed)
	Calibration = 4
	// NormalizerError is used in DSpot when something bad occured during normalization
	NormalizerError = 5
)

// Spot This object embeds a pointer to a C++ object Spot
type Spot struct {
	config *SpotConfig
	status *SpotStatus
	up     *Tail
	down   *Tail
	tmp    []float64
}

// NewSpotFromConfig creates from a SpotConfig structure
func NewSpotFromConfig(conf *SpotConfig) *Spot {
	return &Spot{
		config: conf,
		status: NewSpotStatus(),
		up:     NewTail(conf.MaxExcess),
		down:   NewTail(conf.MaxExcess),
		tmp:    make([]float64, 0),
	}
}

// NewDefaultSpot is the default Spot constructor
func NewDefaultSpot() *Spot {
	return NewSpotFromConfig(&DefaultSpotConfig)
}

func (s *Spot) calibrate() {
	sort.Float64s(s.tmp)

	if s.config.Up {
		// retrieve the upper t threshold
		indexUp := int(s.config.Level * float64(s.config.Ninit))
		s.status.TUp = s.tmp[indexUp-1]

		// feed the tail with the excesses
		for _, ex := range s.tmp[indexUp:] {
			s.status.NtUp++
			s.up.AddExcess(ex - s.status.TUp)
		}
		// upperTail fit
		s.up.Fit()
		s.updateUpThreshold()
	}

	if s.config.Up {
		// retrieve the lower t threshold
		indexDown := int((1. - s.config.Level) * float64(s.config.Ninit))
		s.status.TDown = s.tmp[indexDown]

		// feed the tail with the excesses
		for _, ex := range s.tmp[:indexDown] {
			s.status.NtDown++
			s.down.AddExcess(s.status.TDown - ex)
		}

		// upperTail fit
		s.down.Fit()
		s.updateDownThreshold()
	}
}

func (s *Spot) updateUpThreshold() {
	s.status.ExUp = s.up.ubend.Length()
	s.status.ZUp = s.up.Quantile(
		s.config.Q,
		s.status.TUp,
		s.status.N,
		s.status.NtUp,
	)
}

func (s *Spot) updateDownThreshold() {
	s.status.ExDown = s.down.ubend.Length()
	s.status.ZDown = 2*s.status.TDown - s.down.Quantile(
		s.config.Q,
		s.status.TDown,
		s.status.N,
		s.status.NtDown,
	)
}

// Step performs one Spot step (it analyzes the input data)
func (s *Spot) Step(x float64) int {
	if len(s.tmp) == s.config.Ninit-1 {
		// last init batch data + calibration
		s.tmp = append(s.tmp, x)
		s.status.N++
		s.calibrate()
		return Calibration
	}
	if len(s.tmp) < s.config.Ninit {
		// init batch data
		s.status.N++
		s.tmp = append(s.tmp, x)
		return InitBatch
	}
	if s.config.Up {
		// Up Alert
		if s.config.Alert && x > s.status.ZUp {
			s.status.AlUp++
			return AlertUp
		}
		// Up Excess
		if x > s.status.TUp {
			s.up.AddExcess(x - s.status.TUp)
			s.up.Fit()
			s.updateUpThreshold()
			s.status.NtUp++
			s.status.N++
			return ExcessUp
		}
	}
	if s.config.Down {
		// Down alert
		if s.config.Alert && x < s.status.ZDown {
			s.status.AlDown++
			return AlertDown
		}
		// Down excess
		if x < s.status.TDown {
			s.down.AddExcess(s.status.TDown - x)
			s.down.Fit()
			s.updateDownThreshold()
			s.status.NtDown++
			s.status.N++
			return ExcessDown
		}
	}

	// Normal data
	s.status.N++
	return Normal
}

// GetUpperT Returns the upper threshold t
func (s *Spot) GetUpperT() float64 {
	return s.status.TUp
}

// GetLowerT Returns the lower threshold t
func (s *Spot) GetLowerT() float64 {
	return s.status.TDown
}

// GetUpperThreshold returns the upper decision threshold
func (s *Spot) GetUpperThreshold() float64 {
	return s.status.ZUp
}

// GetLowerThreshold returns the lower decision threshold
func (s *Spot) GetLowerThreshold() float64 {
	return s.status.ZDown
}

// SetQ Change the value of the decision probability.
// It then changes the decision thresholds
func (s *Spot) SetQ(q float64) {
	s.config.Q = q
}

// UpProbability Given a quantile z, computes the probability
// to observe a value greater than z
func (s *Spot) UpProbability(z float64) float64 {
	if s.config.Up {
		return s.up.Cdf(
			z,
			s.status.TUp,
			s.status.N,
			s.status.NtUp)
	}
	return math.NaN()
}

// DownProbability Given a quantile z, computes the probability
// to observe a value lower than z
func (s *Spot) DownProbability(z float64) float64 {
	if s.config.Down {
		return s.down.Cdf(
			2*s.status.TDown-z,
			s.status.TDown,
			s.status.N,
			s.status.NtDown)
	}
	return math.NaN()
}

// Status returns the current status of the Spot instance
func (s *Spot) Status() SpotStatus {
	return *s.status
}

// Config returns the configuration of the Spot instance
func (s *Spot) Config() SpotConfig {
	return *s.config
}
