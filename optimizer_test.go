// optimizer_test.go

package gospot

import (
	"math"
	"testing"
)

func TestInitOptimizer(t *testing.T) {
	title("Optimizers")
}

func parabol(x float64, a interface{}) float64 {
	return 1. + (x-a.(float64))*(x-a.(float64))
}

func fun0(x float64, k interface{}) float64 {
	return -math.Pow(x, k.(float64)) * math.Exp(-x)
}

func TestParabol(t *testing.T) {
	min := 2.0
	a := -10.
	b := 50.
	tol := 1e-8
	xmin, _, _, err := BrentMinimizer(parabol, min, a, b, tol)
	if err != nil {
		t.Fatal(err)
	}
	if (xmin - min) > tol {
		t.Errorf("Minimum not found with given tolerance (expected %f, got %f)", a, xmin)
	}
}

func TestFun0(t *testing.T) {
	k := 7.0
	a := -10.
	b := 200.
	tol := 1e-2

	checkTitle("Brent minimizer...")
	xmin, _, _, err := BrentMinimizer(fun0, k, a, b, tol)
	if err != nil {
		t.Fatal(err)
	}
	if (xmin - k) > tol {
		t.Errorf("Minimum not found with given tolerance (expected %f, got %f)", k, xmin)
		testERROR()
	} else {
		testOK()
	}
}

func TestRoot(t *testing.T) {
	k := 7.0
	a := -10.
	b := 200.
	tol := 1e-8

	checkTitle("Brent root finder...")
	root, _ := BrentRootFinder(fun0, k, a, b, tol)
	if (root - k) > tol {
		t.Errorf("Minimum not found with given tolerance (expected %f, got %f)", k, root)
		testERROR()
	} else {
		testOK()
	}
}

func TestBisection(t *testing.T) {
	k := 7.0
	a := -10.
	b := 200.
	tol := 1e-8

	checkTitle("Bisection...")
	root, _ := Bisection(fun0, k, a, b, tol)
	if (root - k) > tol {
		t.Errorf("Minimum not found with given tolerance (expected %f, got %f)", k, root)
		testERROR()
	} else {
		testOK()
	}
}
