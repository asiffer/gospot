package gospot

import (
	"math"
	"testing"
)

func line(x float64, extra interface{}) float64 {
	return 8*x - 3
}

func square(x float64, extra interface{}) float64 {
	return x*x - 4
}

func noroot(x float64, extra interface{}) float64 {
	return 1.0 + math.Exp(-x)
}

// Example usage
func TestBrent(t *testing.T) {
	tol := 1e-8

	root, found := Brent(0., 3., line, nil, tol)
	if !found {
		t.Errorf("root not found within maximum iterations")
	}
	if math.Abs(root-3.0/8.0) > tol {
		t.Errorf("bad root: %v", root)
	}

	root, found = Brent(0., 3., square, nil, tol)
	if !found {
		t.Errorf("root not found within maximum iterations")
	}
	if math.Abs(root-2.0) > tol {
		t.Errorf("bad root: %v", root)
	}

	root, found = Brent(0., 3., noroot, nil, tol)
	if found {
		t.Errorf("root found while it does not exist")
	}
	if !math.IsNaN(root) {
		t.Errorf("root must be NaN: %v", root)
	}
}
