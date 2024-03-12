package gospot

import (
	"fmt"
	"math"
)

type SpotStatus int

const (
	INTERNAL_ERROR SpotStatus = iota - 1
	NORMAL
	EXCESS
	ANOMALY
)

// Spot represents the main structure to run the SPOT algorithm
type Spot struct {
	// Probability of an anomaly
	Q float64 `json:"q"`
	// Upper/Lower tail choice (1 = lower tail, 0 = upper tail)
	Low bool `json:"low"`
	// Location of the tail (high quantile)
	Level float64 `json:"level"`
	// Flag anomalies (1 = flag, 0 = don't flag)
	DiscardAnomalies bool `json:"discard_anomalies"`
	// Total number of excesses
	Nt uint64 `json:"Nt"`
	// Total number of seen data
	N uint64 `json:"n"`
	// GPD Tail
	Tail *Tail `json:"tail"`
	// Normal/abnormal threshold
	AnomalyThreshold float64 `json:"anomaly_threshold"`
	// Tail threshold
	ExcessThreshold float64 `json:"excess_threshold"`
}

// NewSpot initializes and returns a new Spot instance with the given parameters.
//
// Parameters:
//   - q: Decision probability (Spot will flag extreme events that will have
//     a probability lower than q)
//   - low: Lower tail mode (false for upper tail and true for lower tail)
//   - discardAnomalies: Do not include anomalies in the model (generally true)
//   - level: Excess level (it is a high quantile that delimits the tail)
//   - maxExcess: Maximum number of data that are kept to analyze the tail
//
// Returns:
//   - a pointer to the newly created Spot instance
//   - an error value indicating whether an error occurred during initialization.
//     In particular you must have 0 < level < 1-q < 1
func NewSpot(q float64, low bool, discardAnomalies bool, level float64, maxExcess uint64) (*Spot, error) {
	if level < 0.0 || level >= 1.0 {
		return nil, fmt.Errorf("level must be in [0, 1), close to 1")
	}
	if q >= (1.0-level) || q <= 0.0 {
		return nil, fmt.Errorf("q must be in (0, 1-level)")
	}

	return &Spot{
		Q:                q,
		Level:            level,
		Low:              low,
		DiscardAnomalies: discardAnomalies,
		Nt:               0,
		N:                0,
		Tail:             NewTail(maxExcess),
		AnomalyThreshold: math.NaN(),
		ExcessThreshold:  math.NaN(),
	}, nil
}

func (spot *Spot) upDown() float64 {
	if spot.Low {
		return -1.0
	}
	return 1.0
}

// Reset puts the Spot object into its initial state (before fitting)
func (s *Spot) Reset() {
	maxExcess := uint64(len(s.Tail.Peaks.Container.Data))
	s.Nt = 0
	s.N = 0
	s.Tail = NewTail(maxExcess)
	s.AnomalyThreshold = math.NaN()
	s.ExcessThreshold = math.NaN()
}

// Fit the Spot instance against the given values.
// It computes the excess and anomaly thresholds.
func (spot *Spot) Fit(data []float64) error {
	spot.Nt = 0
	spot.N = uint64(len(data))

	var et float64
	if spot.Low {
		et = P2Quantile(1.0-spot.Level, data)
	} else {
		et = P2Quantile(spot.Level, data)
	}
	if math.IsNaN(et) {
		return fmt.Errorf("excess threshold is NaN")
	}
	spot.ExcessThreshold = et

	for _, x := range data {
		excess := spot.upDown() * (x - et)
		if excess > 0 {
			spot.Nt++
			spot.Tail.Push(excess)
		}
	}

	spot.Tail.Fit()

	spot.AnomalyThreshold = spot.Quantile(spot.Q)
	if math.IsNaN(spot.AnomalyThreshold) {
		return fmt.Errorf("anomaly threshold is NaN")
	}

	return nil
}

// Step updates the Spot instance with a fresh value x
// It returns:
//   - [ANOMALY]: the data is higher the anomaly threshold (or lower in case of lower-tail flagging)
//   - [EXCESS]: the data is in the tail of the distribution and has triggered a model update
//   - [NORMAL]: nothing to say
//   - [INTERNAL_ERROR]: the input value is NaN
func (spot *Spot) Step(x float64) SpotStatus {
	if math.IsNaN(x) {
		return INTERNAL_ERROR
	}

	// flag anomaly
	if spot.DiscardAnomalies && spot.upDown()*(x-spot.AnomalyThreshold) > 0 {
		return ANOMALY
	}

	spot.N++

	ex := spot.upDown() * (x - spot.ExcessThreshold)
	if ex >= 0.0 {
		spot.Nt++
		spot.Tail.Push(ex)
		spot.Tail.Fit()
		spot.AnomalyThreshold = spot.Quantile(spot.Q)
		return EXCESS
	}

	return NORMAL
}

// Quantile computes the value zq such that P(X>zq) = q
func (spot *Spot) Quantile(q float64) float64 {
	s := float64(spot.Nt) / float64(spot.N)
	return spot.ExcessThreshold + spot.upDown()*spot.Tail.Quantile(s, q)
}

// Probability computes the probability p such that P(X>z) = p
func (spot *Spot) Probability(z float64) float64 {
	s := float64(spot.Nt) / float64(spot.N)
	return spot.Tail.Probability(s, spot.upDown()*(z-spot.ExcessThreshold))
}
