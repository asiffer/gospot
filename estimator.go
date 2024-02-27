package gospot

import (
	"math"
)

type Estimator func() (float64, float64, float64)

// MomEstimator computes the 'Method of Moments' estimator for a GPD distribution
func (peaks *Peaks) MomEstimator() (gamma, sigma float64, llhood float64) {
	E := peaks.Mean()
	V := peaks.Var()
	R := E * E / V

	gamma = 0.5 * (1.0 - R)
	sigma = 0.5 * E * (1.0 + R)
	llhood = peaks.LogLikelihood(gamma, sigma)
	return
}

func grimshawW(x float64, extra interface{}) float64 {
	peaks := extra.(*Peaks)
	NtLocal := peaks.Size()
	u := 0.0
	v := 0.0

	for i := uint64(0); i < NtLocal; i++ {
		s := 1.0 + x*peaks.Container.Data[i]
		u += 1 / s
		v += math.Log(s)
	}
	return (u/float64(NtLocal))*(1.0+v/float64(NtLocal)) - 1.0
}

func (peaks *Peaks) grimshawV(x float64) float64 {
	v := 0.0
	NtLocal := peaks.Size()
	for i := 0; i < int(NtLocal); i++ {
		v += math.Log(1.0 + x*peaks.Container.Data[i])
	}
	return 1.0 + v/float64(NtLocal)
}

func (peaks *Peaks) grimshawSimplifiedLogLikelihood(xStar float64) (gamma, sigma, llhood float64) {
	if xStar == 0 {
		gamma = 0.0
		sigma = peaks.Mean()
	} else {
		gamma = peaks.grimshawV(xStar) - 1
		sigma = gamma / xStar
	}
	return gamma, sigma, peaks.LogLikelihood(gamma, sigma)
}

// Grimshaw computes the Grimshaw's estimator for a GPD distribution
func (peaks *Peaks) GrimshawEstimator() (float64, float64, float64) {
	mini := peaks.Min
	maxi := peaks.Max
	mean := peaks.Mean()

	// 0 is always root
	gamma, sigma, maxLLhood := peaks.grimshawSimplifiedLogLikelihood(0.0)

	epsilon := math.Min(BrentDefaultEpsilon, 0.5/maxi)
	a, b := -1.0/maxi+epsilon, -epsilon

	leftRoot, _ := Brent(a, b, grimshawW, peaks, BrentDefaultEpsilon)
	rightRoot, _ := Brent(epsilon, 2.0*(mean-mini)/(mini*mini), grimshawW, peaks, BrentDefaultEpsilon)

	if !math.IsNaN(leftRoot) {
		g, s, ll := peaks.grimshawSimplifiedLogLikelihood(leftRoot)
		if ll > maxLLhood {
			gamma, sigma, maxLLhood = g, s, ll
		}
	}
	if !math.IsNaN(rightRoot) {
		g, s, ll := peaks.grimshawSimplifiedLogLikelihood(rightRoot)
		if ll > maxLLhood {
			gamma, sigma, maxLLhood = g, s, ll
		}
	}

	return gamma, sigma, maxLLhood
}
