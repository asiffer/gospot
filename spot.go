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

// Spot represents a Spot structure
type Spot struct {
	Q                float64 `json:"q"`
	Low              bool    `json:"low"`
	DiscardAnomalies bool    `json:"discard_anomalies"`
	Level            float64 `json:"level"`
	Nt               int     `json:"Nt"`
	N                uint64  `json:"n"`
	Tail             *Tail   `json:"tail"`
	AnomalyThreshold float64 `json:"anomaly_threshold"`
	ExcessThreshold  float64 `json:"excess_threshold"`
}

// NewSpot initializes a new Spot structure
func NewSpot(q float64, low bool, discardAnomalies bool, level float64, maxExcess uint64) (*Spot, error) {
	if level < 0.0 || level >= 1.0 {
		return nil, fmt.Errorf("level must be in [0, 1)")
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

// Fit calculates the excess thresholds
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

// Step updates Spot with a new value
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

// Quantile calculates the quantile for Spot
func (spot *Spot) Quantile(q float64) float64 {
	s := float64(spot.Nt) / float64(spot.N)
	return spot.ExcessThreshold + spot.upDown()*spot.Tail.Quantile(s, q)
}

// Probability calculates the probability for Spot
func (spot *Spot) Probability(z float64) float64 {
	s := float64(spot.Nt) / float64(spot.N)
	return spot.Tail.Probability(s, spot.upDown()*(z-spot.ExcessThreshold))
}
