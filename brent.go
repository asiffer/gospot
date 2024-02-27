package gospot

import (
	"math"
)

const (
	BrentDefaultEpsilon = 2.0e-8
	BrentItmax          = 200
)

func fabs(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}

type RealFunction func(float64, interface{}) float64

func Brent(x1, x2 float64, f RealFunction, extra interface{}, tol float64) (float64, bool) {
	// Assume we found the root
	// It will be set to zero in error cases
	found := true

	a := x1
	b := x2
	c := x2
	d := 0.0
	e := 0.0

	fa := f(a, extra)
	fb := f(b, extra)
	fc := 0.0

	if (fa > 0.0 && fb > 0.0) || (fa < 0.0 && fb < 0.0) {
		return math.NaN(), false
	}

	fc = fb
	for iter := 0; iter < BrentItmax; iter++ {
		if (fb > 0.0 && fc > 0.0) || (fb < 0.0 && fc < 0.0) {
			c = a
			fc = fa
			e = b - a
			d = e
		}
		if fabs(fc) < fabs(fb) {
			a = b
			b = c
			c = a
			fa = fb
			fb = fc
			fc = fa
		}
		tol1 := 2.0*BrentDefaultEpsilon*fabs(b) + 0.5*tol
		xm := 0.5 * (c - b)
		if fabs(xm) <= tol1 || fb == 0.0 {
			return b, found
		}
		if fabs(e) >= tol1 && fabs(fa) > fabs(fb) {
			s := fb / fa
			var p, q float64
			if a == c {
				p = 2.0 * xm * s
				q = 1.0 - s
			} else {
				q = fa / fc
				r := fb / fc
				p = s * (2.0*xm*q*(q-r) - (b-a)*(r-1.0))
				q = (q - 1.0) * (r - 1.0) * (s - 1.0)
			}
			if p > 0.0 {
				q = -q
			}
			p = fabs(p)
			min1 := 3.0*xm*q - fabs(tol1*q)
			min2 := fabs(e * q)
			if 2.0*p < math.Min(min1, min2) {
				e = d
				d = p / q
			} else {
				d = xm
				e = d
			}
		} else {
			d = xm
			e = d
		}
		a = b
		fa = fb
		if fabs(d) > tol1 {
			b += d
		} else {
			if xm >= 0.0 {
				b += fabs(tol1)
			} else {
				b -= fabs(tol1)
			}
		}
		fb = f(b, extra)
	}
	// Maximum number of iterations exceeded
	return math.NaN(), false
}
