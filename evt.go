// evt.go

package gospot

import (
	"math"
)

// Tail defines a distribution tail (EVT framework)
type Tail struct {
	sigma  float64
	gamma  float64
	llhood float64
	ubend  *Ubend
}

var (
	epsilon = 3.e-10
)

func grimshawU(x float64, excesses []float64) float64 {
	u := 0.
	for _, yi := range excesses {
		u += 1. / (1. + x*yi)
	}
	return u / float64(len(excesses))
}

func grimshawV(x float64, excesses []float64) float64 {
	v := 0.
	for _, yi := range excesses {
		v += math.Log(1. + x*yi)
	}
	return 1. + v/float64(len(excesses))
}

func loglikelihoodGPD(sigma float64, gamma float64, excesses []float64) float64 {
	Nt := float64(len(excesses))
	var ll float64
	if gamma == 0. {
		ll = Nt*math.Log(sigma) + sum(excesses)/sigma
	} else {
		a := (1. + 1./gamma)
		b := gamma / sigma
		ll = Nt * math.Log(sigma)
		for _, yi := range excesses {
			ll += a * math.Log(1.+b*yi)
		}
	}
	return -ll
}

// NewTail creates a new tail
func NewTail(size int) *Tail {
	return &Tail{
		gamma:  0.,
		sigma:  0.,
		llhood: 0.,
		ubend:  NewUbend(size)}
}

// AddExcess push a new excess in the tail
func (tail *Tail) AddExcess(x float64) {
	tail.ubend.Push(x)
}

// Fit find the best tail parameters according to
// the input excesses
func (tail *Tail) Fit() {
	excesses := tail.ubend.data
	var g, s float64
	Nt := float64(len(excesses))
	// fmt.Println("Nt =", Nt)
	fun := func(x float64, args interface{}) float64 {
		ex := (args).([]float64)
		return grimshawU(x, ex)*grimshawV(x, ex) - 1.
		// return math.Log(grimshawU(x, ex)) + math.Log(grimshawV(x, ex))
	}
	Ymean := sum(excesses) / Nt
	Ymin := min(excesses)
	Ymax := max(excesses)

	a := (-1. / Ymax) + epsilon
	b := -epsilon
	c := epsilon
	d := 2 * (Ymean - Ymin) / (Ymin * Ymin)

	// check if the function w is convex close to 0
	isConvex := tail.ubend.Var() >= math.Pow(tail.ubend.Mean(), 2.0)
	roots := []float64{0.0}
	var root float64
	var err error

	if isConvex {
		// right root
		root, err = BrentRootFinder(fun, excesses, c, d, 1e-6)
		if err == nil {
			roots = append(roots, root)
		}

		// left root
		root, err = BrentRootFinder(fun, excesses, a, b, 1e-6)
		if err == nil {
			roots = append(roots, root)
		}
	}

	// if err == nil {
	// 	roots = append(roots, root)
	// }

	// if err1 == nil {
	// 	roots = append(roots, root1)
	// }
	// if err2 == nil {
	// 	roots = append(roots, root2)
	// }

	llmax := math.Inf(-1)
	ll := 0.
	for _, x := range roots {
		if math.Abs(x) > epsilon {
			g = grimshawV(x, excesses) - 1.
			s = g / x
			ll = loglikelihoodGPD(s, g, excesses)
		} else {
			g = 0.
			s = Ymean
			ll = -Nt * (1. + math.Log(s))
		}
		// fmt.Println(x, ll)
		if ll > llmax {
			tail.gamma = g
			tail.sigma = s
			llmax = ll
		}
	}
	tail.llhood = llmax
}

// Quantile computes zq such that P(X>zq) = q
func (tail *Tail) Quantile(q float64, t float64, n, Nt int) float64 {
	r := q * float64(n) / float64(Nt)
	if tail.gamma != 0. {
		return t + (tail.sigma/tail.gamma)*(math.Pow(r, -tail.gamma)-1.)
	}
	return t - tail.sigma*math.Log(r)
}

// Cdf computes P(X>zq)
func (tail *Tail) Cdf(zq float64, t float64, n, Nt int) float64 {
	r := float64(Nt) / float64(n)
	if tail.gamma != 0. {
		return r * math.Pow(1.+(tail.sigma/tail.gamma)*(zq-t), -1./tail.gamma)
	}
	return r * math.Exp(-(zq-t)/tail.sigma)
}

func min(v []float64) float64 {
	size := len(v)
	min := 0.
	if size > 0 {
		min = v[0]
		for i := 1; i < size; i++ {
			if v[i] < min {
				min = v[i]
			}
		}
	}
	return min
}

func max(v []float64) float64 {
	size := len(v)
	max := 0.
	if size > 0 {
		max = v[0]
		for i := 1; i < size; i++ {
			if v[i] > max {
				max = v[i]
			}
		}
	}
	return max
}

func sum(v []float64) float64 {
	s := 0.
	for _, x := range v {
		s += x
	}
	return s
}
