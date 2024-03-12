package gospot

import (
	"math"
	"math/rand"
	"sort"
	"testing"
)

func TestSort5(t *testing.T) {
	data := []float64{5, 3, 4, 1, 2}
	sort5(data)

	for i := 1; i <= 5; i++ {
		if data[i-1] != float64(i) {
			t.Errorf("bad sort: %v", data)
		}
	}

	result := sort.Float64Slice(make([]float64, 5))

	for k := 0; k < 50; k++ {
		data = gaussian(5)
		copy(result, data)
		result.Sort()

		sort5(data)
		for i := 0; i < 5; i++ {
			if data[i] != result[i] {
				t.Errorf("bad sort: %v != %v", data, result)
			}
		}
	}
}

func uniform(size uint64) []float64 {
	out := make([]float64, size)
	for i := uint64(0); i < size; i++ {
		out[i] = rand.Float64()
	}
	return out
}

func uniformX(size uint64, scale float64) []float64 {
	out := make([]float64, size)
	for i := uint64(0); i < size; i++ {
		out[i] = scale * rand.Float64()
	}
	return out
}

func gaussian(size uint64) []float64 {
	out := make([]float64, size)
	for i := uint64(0); i < size; i++ {
		out[i] = rand.NormFloat64()
	}
	return out
}

func TestP2QuantileUniform(t *testing.T) {
	var size uint64 = 2000
	tol := math.Sqrt(float64(size))

	for i := 0; i < 100; i++ {
		data := uniform(size)
		for p := 0.1; p < 1.0; p += 0.1 {
			q := P2Quantile(p, data)
			if math.Abs(p-q) > 2./tol {
				t.Errorf("bad quantile: %v != %v", p, q)
			}
		}
	}
}

// Phi is P(X<=x) when X is N(0, 1)
func Phi(x float64) float64 {
	return 0.5 * (1 + math.Erf(x/math.Sqrt2))
}

// func PhiInv(y float64) float64 {
// 	f := func(x float64, extra interface{}) float64 {
// 		return y - Phi(x)
// 	}
// 	x0, _ := Brent(-10, 10, f, nil, BrentDefaultEpsilon)
// 	return x0
// }

func TestP2QuantileNorm(t *testing.T) {
	var size uint64 = 2000
	tol := math.Sqrt(float64(size))

	err := 0

	for i := 0; i < 100; i++ {
		data := gaussian(size)
		for p := 0.1; p < 1.0; p += 0.1 {
			q := P2Quantile(p, data)
			pth := Phi(q)
			if math.Abs(p-pth) > 2./tol {
				// t.Errorf("bad quantile: %v != %v", p, pth)
				err++
			}
		}
	}

	if err > 2 {
		t.Errorf("too many bad quantile computations: %d/100", err)
	}
}

func TestErrorCase(t *testing.T) {
	p2 := NewP2()
	q := p2.quantile([]float64{})
	if !math.IsNaN(q) {
		t.Errorf("output must be NaN, got %v", q)
	}
}
