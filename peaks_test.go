package gospot

import (
	"math"
	"math/rand"
	"sort"
	"testing"
)

func logGaussian(size uint64) []float64 {
	out := make([]float64, size)
	for i := uint64(0); i < size; i++ {
		out[i] = math.Exp(rand.NormFloat64())
	}
	return out
}

func TestPush(t *testing.T) {
	size := uint64(10)
	p := NewPeaks(size)
	x := 2.0
	for i := uint64(0); i < size; i++ {
		p.Push(x)
	}

	if p.Mean() != x {
		t.Errorf("bad mean: %v != %v", p.Mean(), x)
	}
	if p.Var() != 0.0 {
		t.Errorf("bad variance: %v != %v", p.Var(), 0.)
	}
	if p.Min != x {
		t.Errorf("bad minimum: %v != %v", p.Min, x)
	}
	if p.Max != x {
		t.Errorf("bad maximum: %v != %v", p.Max, x)
	}

	for x = 1.0; x <= float64(size); x += 1.0 {
		p.Push(x)
	}
	if p.Min != 1.0 {
		t.Errorf("bad minimum: %v != %v", p.Min, 1.0)
	}
	if p.Max != float64(size) {
		t.Errorf("bad maximum: %v != %v", p.Max, size)
	}
	if p.Mean() != 5.5 {
		t.Errorf("bad mean: %v != %v", p.Mean(), 5.5)
	}
	if p.Var() != 8.25 {
		t.Errorf("bad variance: %v != %v", p.Var(), 8.25)
	}
}

func TestLikelihoodUniform(t *testing.T) {
	size := uint64(1000)
	p := NewPeaks(size)

	u := uniform(size)
	for _, ui := range u {
		p.Push(ui)
	}
	l0 := p.LogLikelihood(0.0, 1.0)
	lpos := p.LogLikelihood(1.0, 1.0)
	lneg := p.LogLikelihood(-1.0, 1.0)
	if (l0 > lneg) || (lpos > lneg) {
		t.Errorf("bad likelihood: %v > %v or %v > %v", l0, lneg, lpos, lneg)
	}
}

func TestLikelihoodGaussian(t *testing.T) {
	size := uint64(1000)
	p := NewPeaks(size)

	initial := sort.Float64Slice(gaussian(200 * size))
	initial.Sort()
	index := uint64(len(initial)) - size
	u := initial[index:]

	for _, ui := range u {
		p.Push(ui - initial[index])
	}

	l0 := p.LogLikelihood(0.0, 1.0)
	lpos := p.LogLikelihood(1.0, 1.0)
	lneg := p.LogLikelihood(-1.0, 1.0)
	if (!math.IsNaN(lneg) && l0 < lneg) || (!math.IsNaN(lpos) && l0 < lpos) {
		t.Errorf("bad likelihood: %v < %v or %v < %v", l0, lneg, l0, lpos)
	}
}

func TestLikelihoodLogGaussian(t *testing.T) {
	size := uint64(1000)
	p := NewPeaks(size)

	initial := sort.Float64Slice(logGaussian(100 * size))
	initial.Sort()
	index := uint64(len(initial)) - size
	u := initial[index:]

	for _, ui := range u {
		p.Push(ui - initial[index])
	}
	l0 := p.LogLikelihood(0.0, 1.0)
	lpos := p.LogLikelihood(1.0, 1.0)
	lneg := p.LogLikelihood(-1.0, 1.0)
	if (lpos < lneg) || (lpos < l0) {
		t.Errorf("bad likelihood: %v < %v or %v < %v", lpos, lneg, lpos, l0)
	}
}
