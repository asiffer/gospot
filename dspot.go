// dspot.go

package gospot

import "math"

// DSpot is built upon Spot. It adds a moving average
// to estimate the local behavior
type DSpot struct {
	normalizer *Normalizer
	Spot
}

// NewDSpot Constructor
func NewDSpot(depth int, q float64, n_init int32, level float64, up bool, down bool, alert bool, bounded bool, max_excess int32) *DSpot {
	up_u8 := bool2uint8(up)
	down_u8 := bool2uint8(down)
	alert_u8 := bool2uint8(alert)
	bounded_u8 := bool2uint8(bounded)
	return &DSpot{normalizer: NewNormalizer(depth, true, false),
		Spot: Spot{ptr: spot_new(q, n_init, level, up_u8, down_u8, alert_u8, bounded_u8, max_excess),
			Up:   up,
			Down: down}}
}

func NewDSpotFromConfig(dsc DSpotConfig) *DSpot {
	return NewDSpot(dsc.Depth, dsc.Q, dsc.N_init, dsc.Level, dsc.Up, dsc.Down, dsc.Alert, dsc.Bounded, dsc.Max_excess)
}

func NewDefaultDSpot() *DSpot {
	return NewDSpotFromConfig(DefaultDSpotConfig)
}

// Average returns the value of the current model
func (ds *DSpot) Average() float64 {
	return ds.normalizer.Average()
}

func (ds *DSpot) Config() DSpotConfig {
	return DSpotConfig{Depth: ds.normalizer.Depth(), SpotConfig: ds.Spot.Config()}
}

func (ds *DSpot) Status() DSpotStatus {
	status := ds.Spot.Status()
	mean := ds.Average()
	if ds.Down {
		status.T_down += mean
		status.Z_down += mean
	}

	if ds.Up {
		status.T_up += mean
		status.Z_up += mean
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
	} else {
		return 5
	}
}

// GetUpperT Returns the upper threshold t
func (ds *DSpot) GetUpperT() float64 {
	if ds.Up {
		return ds.Average() + ds.Spot.GetUpperT()
	} else {
		return math.NaN()
	}
}

// GetUpperT Returns the lower threshold t
func (ds *DSpot) GetLowerT() float64 {
	if ds.Down {
		return ds.Average() + ds.Spot.GetLowerT()
	} else {
		return math.NaN()
	}
}

// GetUpperT Returns the upper decision threshold
func (ds *DSpot) GetUpperThreshold() float64 {
	if ds.Up {
		return ds.Average() + ds.Spot.GetUpperThreshold()
	} else {
		return math.NaN()
	}
}

// GetUpperT Returns the lower decision threshold
func (ds *DSpot) GetLowerThreshold() float64 {
	if ds.Down {
		return ds.Average() + ds.Spot.GetLowerThreshold()
	} else {
		return math.NaN()
	}
}

// UpProbability Given a quantile z, computes the probability
// to observe a value higher than z
func (ds *DSpot) UpProbability(z float64) float64 {
	if ds.Up {
		return ds.Spot.UpProbability(z - ds.Average())
	} else {
		return math.NaN()
	}
}

// DownProbability Given a quantile z, computes the probability
// to observe a value lower than z
func (ds *DSpot) DownProbability(z float64) float64 {
	if ds.Down {
		return ds.Spot.DownProbability(z - ds.Average())
	} else {
		return math.NaN()
	}
}
