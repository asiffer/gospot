package gospot

import (
	"math"
)

// P2 represents the P2 quantile estimator struct
// See aakinshin.net/posts/p2-quantile-estimator/
type P2 struct {
	q  []float64 // Array to store quantiles
	n  []float64 // Array to store indices
	np []float64 // Array to store adjusted indices
	dn []float64 // Array to store adjustment factors
}

func NewP2() *P2 {
	return &P2{
		q:  make([]float64, 5),
		n:  make([]float64, 5),
		np: make([]float64, 5),
		dn: make([]float64, 5),
	}
}

// sort5 sorts the first 5 elements of a slice in ascending order
func sort5(a []float64) {
	if a[1] < a[0] { // Compare 1st and 2nd element #1
		a[0], a[1] = a[1], a[0]
	}
	if a[3] < a[2] { // Compare 3rd and 4th element #2
		a[3], a[2] = a[2], a[3]
	}
	if a[0] < a[2] { // Compare 1st and 3rd element #3
		// run this if 1st element < 3rd element
		a[1], a[2] = a[2], a[1]
		a[2], a[3] = a[3], a[2]
	} else {
		a[1], a[2] = a[2], a[1]
		a[0], a[1] = a[1], a[0]
	}
	// Now 1st, 2nd and 3rd elements are sorted
	// Sort 5th element into 1st, 2nd and 3rd elements
	if a[4] < a[1] { // #4
		if a[4] < a[0] { // #5
			a[4], a[3] = a[3], a[4]
			a[3], a[2] = a[2], a[3]
			a[2], a[1] = a[1], a[2]
			a[1], a[0] = a[0], a[1]
		} else {
			a[4], a[3] = a[3], a[4]
			a[3], a[2] = a[2], a[3]
			a[2], a[1] = a[1], a[2]
		}
	} else {
		if a[4] < a[2] { // #5
			a[4], a[3] = a[3], a[4]
			a[3], a[2] = a[2], a[3]
		} else {
			a[4], a[3] = a[3], a[4]
		}
	}
	// Sort new 5th element into 2nd, 3rd and 4th
	if a[4] < a[2] { // #6
		if a[4] < a[1] { // #7
			a[4], a[3] = a[3], a[4]
			a[3], a[2] = a[2], a[3]
			a[2], a[1] = a[1], a[2]
		} else {
			a[4], a[3] = a[3], a[4]
			a[3], a[2] = a[2], a[3]
		}
	} else {
		if a[4] < a[3] { // #7
			a[4], a[3] = a[3], a[4]
		}
	}
}

// Init initializes the P2 struct with given p value
func (p2 *P2) Init(p float64) {
	for i := 0; i < 5; i++ {
		p2.q[i] = 0.0
		p2.n[i] = float64(i)
		p2.np[i] = 0.0
		p2.dn[i] = 0.0
	}

	// Set initial values based on p
	p2.np[1] = 2 * p
	p2.np[2] = 4 * p
	p2.np[3] = 2 + 2*p
	p2.np[4] = 4

	p2.dn[1] = p / 2
	p2.dn[2] = p
	p2.dn[3] = (p + 1) / 2
	p2.dn[4] = 1
}

// sign returns the sign of a float64 value
func sign(d float64) float64 {
	switch {
	case d > 0:
		return 1.0
	case d < 0:
		return -1.0
	default:
		return 0.0
	}
}

// linear computes the linear interpolation
func (p2 *P2) linear(i int, d int) float64 {
	return p2.q[i] + float64(d)*(p2.q[i+d]-p2.q[i])/(p2.n[i+d]-p2.n[i])
}

// parabolic computes the parabolic interpolation
func (p2 *P2) parabolic(i int, d int) float64 {
	return p2.q[i] + (float64(d)/(p2.n[i+1]-p2.n[i-1]))*((p2.n[i]-p2.n[i-1]+float64(d))*(p2.q[i+1]-p2.q[i])/(p2.n[i+1]-p2.n[i])+(p2.n[i+1]-p2.n[i]-float64(d))*(p2.q[i]-p2.q[i-1])/(p2.n[i]-p2.n[i-1]))
}

// quantile computes the P2 quantile
func (p2 *P2) quantile(x []float64) float64 {
	size := len(x)
	if size < 5 {
		return math.NaN()
	}

	// Initialize the first 5 elements of q with data
	for i := 0; i < 5; i++ {
		p2.q[i] = x[i]
	}

	// Sort the first 5 elements
	sort5(p2.q)

	// Iterate over the remaining elements
	for j := 5; j < size; j++ {
		xj := x[j]
		if xj < p2.q[0] {
			p2.q[0] = xj
		} else if xj > p2.q[4] {
			p2.q[4] = xj
		} else {
			k := 0
			for xj > p2.q[k] {
				k++
			}
			k--

			// Update indices and adjustment factors
			for i := k + 1; i < 5; i++ {
				p2.n[i] += 1.0
			}
			for i := 0; i < 5; i++ {
				p2.np[i] += p2.dn[i]
			}

			// Update quantile markers
			for i := 1; i < 4; i++ {
				d := p2.np[i] - p2.n[i]
				if (d >= 1 && (p2.n[i+1]-p2.n[i]) > 1) || (d <= -1 && (p2.n[i-1]-p2.n[i]) < -1) {
					d = sign(d)
					qp := p2.parabolic(i, int(d))
					if !(p2.q[i-1] < qp && qp < p2.q[i+1]) {
						qp = p2.linear(i, int(d))
					}
					p2.q[i] = qp
					p2.n[i] += d
				}
			}
		}
	}
	return p2.q[2]
}

// P2Quantile computes the P2 quantile of the given data with the specified p value
func P2Quantile(p float64, data []float64) float64 {
	p2 := NewP2()
	p2.Init(p)
	return p2.quantile(data)
}
