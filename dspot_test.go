// dspot_test.go
package gospot

import (
	"fmt"
	"math"
	"testing"
)

func TestInitAndRunDSpot(t *testing.T) {
	title("Testing DSpot initialization and run")
	// init spot object
	// var depth = 50
	// var q = 1e-4
	// var nInit int32 = 2000
	// var level = 0.99
	// up, down, alert, bounded := true, true, true, true
	// var maxExcess int32 = 200

	checkTitle("Building DSpot...")
	// dspot := NewDSpot(
	// 	depth,
	// 	q,
	// 	nInit,
	// 	level,
	// 	up,
	// 	down,
	// 	alert,
	// 	bounded,
	// 	maxExcess)
	dspot := NewDefaultDSpot()
	testOK()

	// data
	var N = 10000
	data := gaussianSample(N)

	checkTitle("Feeding...")
	for i := 0; i < N; i++ {
		dspot.Step(data[i])
	}
	testOK()

	fmt.Println(dspot.Status())

	checkTitle("Deleting...")
	dspot.Delete()
	testOK()
}

func TestDriftComputation(t *testing.T) {
	title("Testing drift computation")
	// init dspot object

	config := DSpotConfig{
		SpotConfig{
			Q:         1e-4,
			Ninit:     2000,
			Level:     0.99,
			Up:        true,
			Down:      true,
			Alert:     true,
			Bounded:   true,
			MaxExcess: 200,
		},
		500,
	}
	dspot := NewDSpotFromConfig(config)

	// data
	var drift = 10.0
	var N = 10000
	data := gaussianSample(N)

	for i := 0; i < N; i++ {
		dspot.Step(data[i] + drift)
	}

	checkTitle("Checking drift...")
	err := math.Abs(dspot.Average()-drift) / drift
	if err > 5. {
		testERROR()
		t.Errorf("Drift: %.3f (expected: %.3f)\n", dspot.Average(), drift)
	} else if err > 2.5 {
		testWARNING()
	} else {
		testOK()
	}

	dspot.Delete()
}

// func TestDSpotStatus(t *testing.T) {
// 	title("Testing DSpot status")
// 	// init spot object
// 	var depth = 500
// 	var q = 1e-3
// 	var nInit int32 = 2000
// 	var level = 0.99
// 	up, down, alert, bounded := true, true, false, true
// 	var maxExcess int32 = 200

// 	dspot := NewDSpot(
// 		depth,
// 		q,
// 		nInit,
// 		level,
// 		up,
// 		down,
// 		alert,
// 		bounded,
// 		maxExcess)

// 	// data
// 	var drift = -7.0
// 	var N = 12000
// 	data := gaussianSample(N)

// 	for i := 0; i < N; i++ {
// 		dspot.Step(data[i] + drift)
// 	}

// 	fmt.Println(dspot.Status())
// 	dspot.Delete()
// }

func TestNullDepth(t *testing.T) {
	title("Testing null depth")
	// init spot object
	var depth = 0
	var q = 1e-3
	var nInit int32 = 2000
	var level = 0.99
	up, down, alert, bounded := true, true, true, true
	var maxExcess int32 = 200

	dspot := NewDSpot(
		depth,
		q,
		nInit,
		level,
		up,
		down,
		alert,
		bounded,
		maxExcess)
	spot := NewSpot(
		q,
		nInit,
		level,
		up,
		down,
		alert,
		bounded,
		maxExcess)

	// data
	var N = 7000
	data := gaussianSample(N)

	for i := 0; i < N; i++ {
		dspot.Step(data[i])
		spot.Step(data[i])
	}

	checkTitle("Checking Spot/DSpot status...")
	if dspot.Status().SpotStatus != spot.Status() {
		t.Error("Different status")
		testERROR()
		fmt.Print("\n-- DSPOT --\n", dspot.Status(), "\n")
		fmt.Print("-- SPOT --\n", spot.Status(), "\n")
	} else {
		testOK()
	}

	spot.Delete()
	dspot.Delete()
}

func TestDBasicDSpotAccess(t *testing.T) {
	title("Testing DSpot status accesses")

	config := DSpotConfig{
		SpotConfig{
			Q:         1e-4,
			Ninit:     2000,
			Level:     0.99,
			Up:        true,
			Down:      true,
			Alert:     true,
			Bounded:   true,
			MaxExcess: 200,
		},
		500,
	}
	dspot := NewDSpotFromConfig(config)

	N := 2 * int(config.Ninit)

	data := gaussianSample(N)

	for i := 0; i < N; i++ {
		dspot.Step(data[i])
	}

	checkTitle("Checking TUp/TDown...")
	errTUp := math.Abs(dspot.GetUpperT()-dspot.Average()-2.326) / 2.326
	errTDown := math.Abs(dspot.GetLowerT()-dspot.Average()+2.326) / 2.326
	if errTUp > 5 || errTDown > 5 {
		t.Error("Error on the TUp/TDown values")
		fmt.Println(dspot.Status())
		testERROR()
	} else if errTUp > 2.5 || errTDown > 2.5 {
		testWARNING()
	} else {
		testOK()
	}

	checkTitle("Checking ZUp/ZDown...")
	errZUp := math.Abs(dspot.GetUpperThreshold()-dspot.Average()-3.719) / 3.719
	errZDown := math.Abs(dspot.GetLowerThreshold()-dspot.Average()+3.719) / 3.719
	if errZUp > 5 || errZDown > 5 {
		t.Error("Error on the TUp/TDown values")
		fmt.Println("\n", dspot.Status())
		testERROR()
	} else if errZUp > 2.5 || errZDown > 2.5 {
		testWARNING()
	} else {
		testOK()
	}
}

func TestDSpotProbabilityComputation(t *testing.T) {
	title("Probability computation")

	config := DSpotConfig{
		SpotConfig{
			Q:         1e-4,
			Ninit:     2000,
			Level:     0.99,
			Up:        true,
			Down:      true,
			Alert:     true,
			Bounded:   true,
			MaxExcess: 200,
		},
		500,
	}
	dspot := NewDSpotFromConfig(config)
	N := int(config.Ninit)
	data := gaussianSample(N)
	drift := 7.0

	for i := 0; i < N; i++ {
		dspot.Step(data[i] + drift)
	}

	checkTitle("Checking Up probability computation...")
	errUp := math.Abs(dspot.UpProbability(3.09+dspot.Average())-1e-3) / 1e-3
	if errUp > 5 {
		testERROR()
		t.Errorf("Expected 1e-3, got %f", dspot.UpProbability(3.09+dspot.Average()))
	} else if errUp > 2.5 {
		testWARNING()
	} else {
		testOK()
	}

	checkTitle("Checking Down probability computation...")
	errDown := math.Abs(dspot.DownProbability(-3.09+dspot.Average())-1e-3) / 1e-3
	if errDown > 5 {
		testERROR()
		t.Errorf("Expected 1e-3, got %f", dspot.DownProbability(-3.09+dspot.Average()))
	} else if errDown > 2.5 {
		testWARNING()
	} else {
		testOK()
	}

	config = DSpotConfig{
		SpotConfig{
			Q:         1e-4,
			Ninit:     2000,
			Level:     0.99,
			Up:        false,
			Down:      false,
			Alert:     true,
			Bounded:   true,
			MaxExcess: 200,
		},
		500,
	}
	dspot = NewDSpotFromConfig(config)
	checkTitle("Checking NaN (Up)...")
	checkTitle("Checking NaN (Up)...")
	if math.IsNaN(dspot.UpProbability(12.)) && math.IsNaN(dspot.GetUpperT()) && math.IsNaN(dspot.GetUpperThreshold()) {
		testOK()
	} else {
		testERROR()
	}

	checkTitle("Checking NaN (Down)...")
	if math.IsNaN(dspot.DownProbability(12.)) && math.IsNaN(dspot.GetLowerT()) && math.IsNaN(dspot.GetLowerThreshold()) {
		testOK()
	} else {
		testERROR()
	}

	fmt.Println(dspot.Config())
}
