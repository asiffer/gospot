// dspot.go

package gospot

import (
	"math"
)

// DSpot is built upon Spot. It adds a moving average
// to estimate the local behavior
type DSpot struct {
	normalizer *Normalizer
	Spot
}

// NewDSpotFromConfig creates a DSpot instance from a config structure
func NewDSpotFromConfig(dsc *DSpotConfig) *DSpot {
	return &DSpot{
		normalizer: NewNormalizer(dsc.Depth, true, false),
		Spot: Spot{
			config: &dsc.SpotConfig,
			status: NewSpotStatus(),
			up:     NewTail(dsc.MaxExcess),
			down:   NewTail(dsc.MaxExcess),
			tmp:    make([]float64, 0, dsc.Ninit),
		},
	}
}

// NewDefaultDSpot is the default DSpot constructor
func NewDefaultDSpot() *DSpot {
	return NewDSpotFromConfig(&DefaultDSpotConfig)
}

// Average returns the value of the current model
func (ds *DSpot) Average() float64 {
	return ds.normalizer.Average()
}

// Config returns the initial config of the DSpot instance
func (ds *DSpot) Config() DSpotConfig {
	return DSpotConfig{
		Depth:      ds.normalizer.Depth(),
		SpotConfig: ds.Spot.Config()}
}

// Status returns the current status of the DSpot instance
func (ds *DSpot) Status() DSpotStatus {
	status := ds.Spot.Status()
	mean := ds.Average()
	if ds.config.Down {
		status.TDown += mean
		status.ZDown += mean
	}

	if ds.config.Up {
		status.TUp += mean
		status.ZUp += mean
	}

	return DSpotStatus{Mean: mean, SpotStatus: status}
}

// Step Method which update the Spot instance according to a new incoming value
func (ds *DSpot) Step(x float64) int {
	z, err := ds.normalizer.Step(x)
	if err != nil {
		return NormalizerError
	}
	// normal spot step
	ret := ds.Spot.Step(z)
	if ret == AlertUp || ret == AlertDown {
		// if anomaly, it is not taken in the model
		// cancel the previous Step()
		ds.normalizer.Cancel()
	}
	return ret
}

// GetUpperT Returns the upper threshold t
func (ds *DSpot) GetUpperT() float64 {
	if ds.config.Up {
		return ds.Average() + ds.Spot.GetUpperT()
	}
	return math.NaN()

}

// GetLowerT returns the lower threshold t
func (ds *DSpot) GetLowerT() float64 {
	if ds.config.Down {
		return ds.Average() + ds.Spot.GetLowerT()
	}
	return math.NaN()
}

// GetUpperThreshold returns the upper decision threshold
func (ds *DSpot) GetUpperThreshold() float64 {
	if ds.config.Up {
		return ds.Average() + ds.Spot.GetUpperThreshold()
	}
	return math.NaN()
}

// GetLowerThreshold returns the lower decision threshold
func (ds *DSpot) GetLowerThreshold() float64 {
	if ds.config.Down {
		return ds.Average() + ds.Spot.GetLowerThreshold()
	}
	return math.NaN()
}

// UpProbability Given a quantile z, computes the probability
// to observe a value higher than z
func (ds *DSpot) UpProbability(z float64) float64 {
	if ds.config.Up {
		return ds.Spot.UpProbability(z - ds.Average())
	}
	return math.NaN()
}

// DownProbability Given a quantile z, computes the probability
// to observe a value lower than z
func (ds *DSpot) DownProbability(z float64) float64 {
	if ds.config.Down {
		return ds.Spot.DownProbability(z - ds.Average())
	}
	return math.NaN()
}
