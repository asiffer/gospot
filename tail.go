package gospot

import (
	"math"
)

type Tail struct {
	Gamma float64 `json:"gamma"`
	Sigma float64 `json:"sigma"`
	Peaks *Peaks  `json:"peaks"`
}

func NewTail(size uint64) *Tail {
	return &Tail{
		Gamma: 0.0,
		Sigma: 0.0,
		Peaks: NewPeaks(size),
	}
}

func (tail *Tail) Push(x float64) {
	tail.Peaks.Push(x)
}

func (tail *Tail) Probability(s, d float64) float64 {
	// d = zq - t
	if tail.Gamma == 0.0 {
		return s * math.Exp(-d/tail.Sigma)
	} else {
		r := d * (tail.Gamma / tail.Sigma)
		return s * math.Pow(1.0+r, -1.0/tail.Gamma)
	}
}

func (tail *Tail) Quantile(s, q float64) float64 {
	r := q / s
	if tail.Gamma == 0.0 {
		return -tail.Sigma * math.Log(r)
	}
	return (tail.Sigma / tail.Gamma) * (math.Pow(r, -tail.Gamma) - 1)
}

func (tail *Tail) Fit() float64 {
	maxLLhood := math.NaN()

	for _, e := range []Estimator{tail.Peaks.MomEstimator, tail.Peaks.GrimshawEstimator} {
		gamma, sigma, llhood := e()
		if math.IsNaN(maxLLhood) || (llhood > maxLLhood) {
			maxLLhood = llhood
			tail.Gamma = gamma
			tail.Sigma = sigma
		}
	}

	return maxLLhood
}
