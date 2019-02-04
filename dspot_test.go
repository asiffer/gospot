// dspot_test.go
package gospot

import (
	"fmt"
	"strings"
	"testing"
)

var (
	HEADER_WIDTH int    = 100
	HEADER_SYM   string = "="
)

func title(s string) {
	var l int = len(s)
	var border int
	var left string
	var right string
	remaining := HEADER_WIDTH - l - 2
	if remaining%2 == 0 {
		border = remaining / 2
		left = strings.Repeat(HEADER_SYM, border) + " "
		right = " " + strings.Repeat(HEADER_SYM, border)
	} else {
		border = (remaining - 1) / 2
		left = strings.Repeat(HEADER_SYM, border+1) + " "
		right = " " + strings.Repeat(HEADER_SYM, border)
	}

	fmt.Println(left + s + right)
}

func TestDSpot(t *testing.T) {
	title("DSPOT")
}

func TestInitAndRunDSpot(t *testing.T) {
	fmt.Println("** Testing initialization and run **")
	// init spot object
	var depth int = 50
	var q float64 = 1e-4
	var n_init int32 = 2000
	var level float64 = 0.99
	up, down, alert, bounded := true, true, true, true
	var max_excess int32 = 200

	dspot := NewDSpot(depth, q, n_init, level, up, down, alert, bounded, max_excess)
	// data
	var N int = 10000
	data := gaussianSample(N)

	for i := 0; i < N; i++ {
		dspot.Step(data[i])
	}

	dspot.Delete()
}

func TestDriftComputation(t *testing.T) {
	fmt.Println("** Testing drift computation **")
	// init spot object
	var depth int = 500
	var q float64 = 1e-4
	var n_init int32 = 2000
	var level float64 = 0.99
	up, down, alert, bounded := true, true, true, true
	var max_excess int32 = 200

	dspot := NewDSpot(depth, q, n_init, level, up, down, alert, bounded, max_excess)

	// data
	var drift float64 = 10.0
	var N int = 10000
	data := gaussianSample(N)

	for i := 0; i < N; i++ {
		dspot.Step(data[i] + drift)
	}

	fmt.Printf("Drift: %.3f (expected: %.3f)\n", dspot.Average(), drift)
	dspot.Delete()
}

func TestDSpotStatus(t *testing.T) {
	fmt.Println("** Testing status **")
	// init spot object
	var depth int = 500
	var q float64 = 1e-3
	var n_init int32 = 2000
	var level float64 = 0.99
	up, down, alert, bounded := true, true, true, true
	var max_excess int32 = 200

	dspot := NewDSpot(depth, q, n_init, level, up, down, alert, bounded, max_excess)

	// data
	var drift float64 = -7.0
	var N int = 12000
	data := gaussianSample(N)

	for i := 0; i < N; i++ {
		dspot.Step(data[i] + drift)
	}

	fmt.Println(dspot.Status())
	dspot.Delete()
}

func TestNullDepth(t *testing.T) {
	fmt.Println("** Testing null depth **")
	// init spot object
	var depth int = 0
	var q float64 = 1e-3
	var n_init int32 = 2000
	var level float64 = 0.99
	up, down, alert, bounded := true, true, true, true
	var max_excess int32 = 200

	dspot := NewDSpot(depth, q, n_init, level, up, down, alert, bounded, max_excess)
	spot := NewSpot(q, n_init, level, up, down, alert, bounded, max_excess)
	// data

	var N int = 7000
	data := gaussianSample(N)

	for i := 0; i < N; i++ {
		dspot.Step(data[i])
		spot.Step(data[i])
	}

	fmt.Print("-- DSPOT --\n", dspot.Status(), "\n")
	fmt.Print("-- SPOT --\n", spot.Status(), "\n")

	if dspot.Status().SpotStatus != spot.Status() {
		t.Error("Different status")
	}

	spot.Delete()
	dspot.Delete()
}
