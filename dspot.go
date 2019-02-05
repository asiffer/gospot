// dspot.go

package gospot

import "math"

// DSpot is built upon Spot. It adds a moving average
// to estimate the local behavior
type DSpot struct {
	normalizer *Normalizer
	Spot
}

// NewDSpot is the DSpot constructor
func NewDSpot(depth int, q float64, nInit int32,
	level float64, up bool, down bool,
	alert bool, bounded bool, maxExcess int32) *DSpot {

	up8 := bool2uint8(up)
	down8 := bool2uint8(down)
	alert8 := bool2uint8(alert)
	bounded8 := bool2uint8(bounded)

	return &DSpot{normalizer: NewNormalizer(depth, true, false),
		Spot: Spot{
			ptr: spotNew(q, nInit, level, up8,
				down8, alert8, bounded8, maxExcess),
			Up:   up,
			Down: down}}
}

// NewDSpotFromConfig creates a DSpot instance from a config structure
func NewDSpotFromConfig(dsc DSpotConfig) *DSpot {
	return NewDSpot(dsc.Depth, dsc.Q, dsc.Ninit, dsc.Level, dsc.Up, dsc.Down, dsc.Alert, dsc.Bounded, dsc.MaxExcess)
}

// NewDefaultDSpot is the default DSpot constructor
func NewDefaultDSpot() *DSpot {
	return NewDSpotFromConfig(DefaultDSpotConfig)
}

// Average returns the value of the current model
func (ds *DSpot) Average() float64 {
	return ds.normalizer.Average()
}

// Config returns the initial config of the DSpot instance
func (ds *DSpot) Config() DSpotConfig {
	return DSpotConfig{Depth: ds.normalizer.Depth(), SpotConfig: ds.Spot.Config()}
}

// Status returns the current status of the DSpot instance
func (ds *DSpot) Status() DSpotStatus {
	status := ds.Spot.Status()
	mean := ds.Average()
	if ds.Down {
		status.TDown += mean
		status.ZDown += mean
	}

	if ds.Up {
		status.TUp += mean
		status.ZUp += mean
	}

	return DSpotStatus{Mean: mean, SpotStatus: status}
}

// Step Method which update the Spot instance according to a new incoming value
func (ds *DSpot) Step(x float64) int32 {
	z, err := ds.normalizer.Step(x)
	if err == nil {
		ret := ds.Spot.Step(z)
		if ret*ret == 1 { // if anomaly, it is not taken in the model
			// cancel the previous Step()
			ds.normalizer.Cancel()
		}
		return ret
	}
	// error code
	return 5
}

// GetUpperT Returns the upper threshold t
func (ds *DSpot) GetUpperT() float64 {
	if ds.Up {
		return ds.Average() + ds.Spot.GetUpperT()
	}
	return math.NaN()

}

// GetLowerT returns the lower threshold t
func (ds *DSpot) GetLowerT() float64 {
	if ds.Down {
		return ds.Average() + ds.Spot.GetLowerT()
	}
	return math.NaN()
}

// GetUpperThreshold returns the upper decision threshold
func (ds *DSpot) GetUpperThreshold() float64 {
	if ds.Up {
		return ds.Average() + ds.Spot.GetUpperThreshold()
	}
	return math.NaN()
}

// GetLowerThreshold returns the lower decision threshold
func (ds *DSpot) GetLowerThreshold() float64 {
	if ds.Down {
		return ds.Average() + ds.Spot.GetLowerThreshold()
	}
	return math.NaN()
}

// UpProbability Given a quantile z, computes the probability
// to observe a value higher than z
func (ds *DSpot) UpProbability(z float64) float64 {
	if ds.Up {
		return ds.Spot.UpProbability(z - ds.Average())
	}
	return math.NaN()
}

// DownProbability Given a quantile z, computes the probability
// to observe a value lower than z
func (ds *DSpot) DownProbability(z float64) float64 {
	if ds.Down {
		return ds.Spot.DownProbability(z - ds.Average())
	}
	return math.NaN()
}
