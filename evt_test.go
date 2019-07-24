// evt_test.go

package gospot

import (
	"math"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestEVT(t *testing.T) {
	title("EVT Tail Fit")
	rand.Seed(time.Now().UnixNano())
}

func TestCdf(t *testing.T) {
	N := 10000
	Nt := N / 100

	checkTitle("Checking tail fit (normal)...")
	tail := NewTail(-1)

	data := standardGaussianSample(N)
	sort.Float64s(data)

	for i := N - Nt; i < N; i++ {
		tail.AddExcess(data[i] - data[N-Nt-1])
	}
	tail.Fit()

	if tail.gamma > 0. {
		t.Errorf("Bad fitted gamma, expected 0., got %f", tail.gamma)
		testERROR()
	} else if tail.gamma < 0 {
		testWARNING()
		t.Errorf("Bad fitted gamma, expected 0., got %f", tail.gamma)
	} else {
		testOK()
	}

	checkTitle("Checking tail fit (uniform)...")
	tail = NewTail(-1)
	data = uniformSample(N)
	sort.Float64s(data)

	for i := N - Nt; i < N; i++ {
		tail.AddExcess(data[i] - data[N-Nt-1])
	}
	tail.Fit()

	if tail.gamma > 0. {
		t.Errorf("Bad fitted gamma, expected γ<0, got %f", tail.gamma)
		testERROR()
	} else if tail.gamma == 0 {
		testWARNING()
	} else {
		testOK()
	}

}

func TestQuantile(t *testing.T) {
	N := 50000
	Nt := 100
	checkTitle("Checking quantile (normal)...")
	tail := NewTail(-1)

	data := standardGaussianSample(N)
	sort.Float64s(data)
	for i := N - Nt; i < N; i++ {
		tail.AddExcess(data[i] - data[N-Nt-1])
	}
	tail.Fit()

	q := 1e-4
	zq := tail.Quantile(q, data[N-Nt-1], N, Nt)
	if math.Abs(3.72-zq)/3.72 > 0.05 {
		t.Errorf("Bad quantile computation, expected 3.72, got %f", zq)
		testERROR()
	} else {
		testOK()
	}

	checkTitle("Checking quantile (uniform)...")
	tail = NewTail(-1)
	data = uniformSample(N)
	sort.Float64s(data)
	for i := N - Nt; i < N; i++ {
		tail.AddExcess(data[i] - data[N-Nt-1])
	}
	tail.Fit()

	q = 1e-5
	zq = tail.Quantile(q, data[N-Nt-1], N, Nt)
	if math.Abs(1-q-zq) > 0.05 {
		t.Errorf("Bad quantile computation, expected 0.9999, got %f", zq)
		testERROR()
	} else {
		testOK()
	}
}
