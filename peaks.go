package gospot

import (
	"math"
)

// Peaks is a stucture that computes stats about the provided excesses
type Peaks struct {
	// Sum of the elements
	E float64 `json:"e"`
	// Sum of the square of the elements
	E2 float64 `json:"e2"`
	// Minimum of the elements
	Min float64 `json:"min"`
	// Maximum of the elements
	Max float64 `json:"max"`
	// Underlying data container
	Container *Ubend `json:"container"`
}

// NewPeaks initializes a new [Peak] structure
func NewPeaks(size uint64) *Peaks {
	return &Peaks{
		E:         0.0,
		E2:        0.0,
		Min:       math.NaN(),
		Max:       math.NaN(),
		Container: NewUbend(size),
	}
}

func (peaks *Peaks) updateStats() uint64 {
	maxIteration := peaks.Container.Size()

	peaks.Min = math.NaN()
	peaks.Max = math.NaN()
	peaks.E = 0.0
	peaks.E2 = 0.0

	for i := uint64(0); i < maxIteration; i++ {
		value := peaks.Container.Data[i]
		peaks.E += value
		peaks.E2 += value * value
		if math.IsNaN(peaks.Min) || (value < peaks.Min) {
			peaks.Min = value
		}
		if math.IsNaN(peaks.Max) || (value > peaks.Max) {
			peaks.Max = value
		}
	}

	return maxIteration
}

// Size returns the current number of peaks
func (peaks *Peaks) Size() uint64 {
	return peaks.Container.Size()
}

// Push a new data to the peaks
func (peaks *Peaks) Push(x float64) {
	erased := peaks.Container.Push(x)
	size := peaks.Size()

	peaks.E += x
	peaks.E2 += x * x

	if (size == 1) || (x < peaks.Min) {
		peaks.Min = x
	}
	if (size == 1) || (x > peaks.Max) {
		peaks.Max = x
	}

	if !math.IsNaN(erased) {
		peaks.E -= erased
		peaks.E2 -= erased * erased
		if (erased <= peaks.Min) || (erased >= peaks.Max) {
			peaks.updateStats()
		}
	}
}

// Mean computes the mean of the peaks
func (peaks *Peaks) Mean() float64 {
	return peaks.E / float64(peaks.Size())
}

// Var computes the variance of the peaks
func (peaks *Peaks) Var() float64 {
	size := float64(peaks.Size())
	mean := peaks.E / size
	return (peaks.E2 / size) - (mean * mean)
}

// LogLikelihood computes the log-likelihood of the peaks against a GPD(gamma, sigma) distribution
func (peaks *Peaks) LogLikelihood(gamma, sigma float64) float64 {
	NtLocal := peaks.Size()
	Nt := float64(NtLocal)

	if gamma == 0.0 {
		return -Nt*math.Log(sigma) - peaks.E/sigma
	}

	r := -Nt * math.Log(sigma)
	c := 1.0 + 1.0/gamma
	x := gamma / sigma

	for i := uint64(0); i < NtLocal; i++ {
		r += -c * math.Log(1+x*peaks.Container.Data[i])
	}

	return r
}
