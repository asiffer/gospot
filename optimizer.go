// optimizer.go

package gospot

import (
	"errors"
	"fmt"
	"math"
)

var (
	// Eps is the machine floating-point precision
	Eps = 3.e-8
	// MaxFunEval is the maximum allowed number of iterations
	MaxFunEval = 200
)

// ObjectiveFunction defines a scalar function to minimize
type ObjectiveFunction func(x float64, args interface{}) float64

// BrentMinimizer minimizes the function f according to the Brent's method
func BrentMinimizer(f ObjectiveFunction, args interface{}, a, b, t float64) (float64, float64, int, error) {
	fEvals := 0

	var c, d, e, eps float64
	var fu, fv, fw, fx float64
	var m, p, q, r float64
	var sa, sb float64
	var t2, tol float64
	var u, v, w, x float64

	//
	//  C is the square of the inverse of the golden ratio.
	//
	c = 0.5 * (3.0 - math.Sqrt(5.0))

	eps = math.Sqrt(2.220446049250313e-16)

	sa = a
	sb = b
	x = sa + c*(b-a)
	w = x
	v = w
	e = 0.0
	fx = f(x, args)
	fEvals++
	fw = fx
	fv = fw

	for fEvals < MaxFunEval {
		m = 0.5 * (sa + sb)
		tol = eps*math.Abs(x) + t
		t2 = 2.0 * tol
		//
		//  Check the stopping criterion.
		//
		if math.Abs(x-m) <= t2-0.5*(sb-sa) {
			return x, fx, fEvals, nil
		}
		//
		//  Fit a parabola.
		//
		r = 0.0
		q = r
		p = q

		if tol < math.Abs(e) {
			r = (x - w) * (fx - fv)
			q = (x - v) * (fx - fw)
			p = (x-v)*q - (x-w)*r
			q = 2.0 * (q - r)
			if 0.0 < q {
				p = -p
			}
			q = math.Abs(q)
			r = e
			e = d
		}

		if math.Abs(p) < math.Abs(0.5*q*r) &&
			q*(sa-x) < p &&
			p < q*(sb-x) {
			//
			//  Take the parabolic interpolation step.
			//
			d = p / q
			u = x + d
			//
			//  F must not be evaluated too close to A or B.
			//
			if (u-sa) < t2 || (sb-u) < t2 {
				if x < m {
					d = tol
				} else {
					d = -tol
				}
			}
		} else {
			//
			//  A golden-section step.
			//
			if x < m {
				e = sb - x
			} else {
				e = sa - x
			}
			d = c * e
		}
		//
		//  F must not be evaluated too close to X.
		//
		if tol <= math.Abs(d) {
			u = x + d
		} else if 0.0 < d {
			u = x + tol
		} else {
			u = x - tol
		}

		fu = f(u, args)
		fEvals++
		//
		//  Update A, B, V, W, and X.
		//
		if fu <= fx {
			if u < x {
				sb = x
			} else {
				sa = x
			}
			v = w
			fv = fw
			w = x
			fw = fx
			x = u
			fx = fu
		} else {
			if u < x {
				sa = u
			} else {
				sb = u
			}

			if fu <= fw || w == x {
				v = w
				fv = fw
				w = u
				fw = fu
			} else if fu <= fv || v == x || v == w {
				v = u
				fv = fu
			}
		}
	}

	return x, fx, fEvals, fmt.Errorf("Maximum number of function evaluations reached")
}

// BrentRootFinder finds a root of the the function f: x->f(x, args)
// between x1 and x2 with the Van Wijngaarden–Dekker–Brent method.
// The implementation directly comes from the book 'Numerical Recipes in C'
// (p. 361, 362)
func BrentRootFinder(f ObjectiveFunction, args interface{}, x1, x2, tol float64) (float64, error) {
	Eps := 3.e-8
	fEvals := 0

	var d, e, min1, min2 float64
	a := x1
	b := x2
	c := x2

	fa := f(a, args)
	fb := f(b, args)
	fEvals += 2
	var fc, p, q, r, s, tol1, xm float64

	if (fa > 0.0 && fb > 0.0) || (fa < 0.0 && fb < 0.0) {
		return 0., errors.New("Root must be bracketed in brent")
	}

	fc = fb
	for fEvals < MaxFunEval {
		if (fb > 0.0 && fc > 0.0) || (fb < 0.0 && fc < 0.0) {
			//  Rename a, b, c and adjust bounding interval
			c = a
			fc = fa
			e = b - a
			d = b - a
		}

		if math.Abs(fc) < math.Abs(fb) {
			a = b
			b = c
			c = a
			fa = fb
			fb = fc
			fc = fa
		}
		//  Convergence check.
		tol1 = 2.0*Eps*math.Abs(b) + 0.5*tol
		xm = 0.5 * (c - b)
		if math.Abs(xm) <= tol1 || fb == 0.0 {
			return b, nil
		}
		if math.Abs(e) >= tol1 && math.Abs(fa) > math.Abs(fb) {
			s = fb / fa
			//  Attempt inverse quadratic interpolation.
			if a == c {
				p = 2.0 * xm * s
				q = 1.0 - s
			} else {
				q = fa / fc
				r = fb / fc
				p = s * (2.0*xm*q*(q-r) - (b-a)*(r-1.0))
				q = (q - 1.0) * (r - 1.0) * (s - 1.0)
			}

			// Check whether in bounds.
			if p > 0.0 {
				q = -q
			}

			p = math.Abs(p)
			min1 = 3.0*xm*q - math.Abs(tol1*q)
			min2 = math.Abs(e * q)

			if 2.0*p < math.Min(min1, min2) {
				// Accept interpolation.
				e = d
				d = p / q
			} else {
				//  Interpolation failed, use bisection.
				d = xm
				e = d
			}
		} else {
			//  Bounds decreasing too slowly, use bisection.
			d = xm
			e = d
		}
		// Move last best guess to a.
		a = b
		fa = fb
		if math.Abs(d) > tol1 {
			// Evaluate new trial root.
			b += d
		} else {
			if xm > 0. {
				b += tol1
			} else {
				b += -tol1
			}
		}
		fb = f(b, args)
		fEvals++
	}
	return 0., fmt.Errorf("Maximum number of function evaluations reached")
}

// Bisection finds a root without derivatives
func Bisection(f ObjectiveFunction, args interface{}, x1, x2, tol float64) (float64, error) {
	fEvals := 0

	a := math.Min(x1, x2)
	b := math.Max(x1, x2)

	fa := f(a, args)
	fb := f(b, args)
	fEvals += 2

	var m, fm float64
	if fa*fb > 0 {
		return 0., errors.New("Root must be bracketed in bisection")
	}

	for (b-a) > tol && fEvals < MaxFunEval {
		m = (a + b) / 2.

		fm = f(m, args)
		fEvals++

		if fa*fm <= 0 {
			b = m
			fb = fm
		} else {
			a = m
			fa = fm
		}

	}

	if fEvals >= MaxFunEval {
		return b, fmt.Errorf("Maximum number of function evaluations reached")
	}
	return b, nil
}

// BFGS uses the BFGS algorithm to find the minimum of a function
// func BFGS(f ObjectiveFunction, args interface{}, x0 float64) (float64, float64, int, error) {
// 	p := optimize.Problem{
// 		Func: func(x []float64) float64 {
// 			return f(x[0], args)
// 		},
// 	}
// 	s := optimize.Settings{
// 		FuncEvaluations: MaxFunEval,
// 	}
// 	result, err := optimize.Minimize(p, []float64{x0}, &s, nil)
// 	return result.X[0], result.F, result.Stats.FuncEvaluations, err
// }
