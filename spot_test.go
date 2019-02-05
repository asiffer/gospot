// spot_test.go

package gospot

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func gaussianSample(N int) []float64 {
	rand.Seed(time.Now().UTC().UnixNano())
	data := make([]float64, N)
	for i := 0; i < N; i++ {
		data[i] = rand.NormFloat64()
	}
	return data
}

func TestSpotInitAndRun(t *testing.T) {
	title("Testing Spot initialization and run")
	// init spot object

	checkTitle("Default Spot creation...")
	spot := NewDefaultSpot()
	testOK()
	// data
	var N = 10000
	data := gaussianSample(N)

	checkTitle("Feeding...")
	for i := 0; i < N; i++ {
		spot.Step(data[i])
	}
	testOK()

	checkTitle("Deleting...")
	spot.Delete()
	testOK()
}

func TestSpotThresholdComputation(t *testing.T) {
	title("Testing Spot threshold computation")
	// init spot object
	var q = 1e-4
	var nInit int32 = 3000
	var level = 0.99
	up, down, alert, bounded := true, true, false, true
	var maxExcess int32 = 300

	spot := NewSpot(q, nInit, level, up, down, alert, bounded, maxExcess)

	checkTitle("Checking Q setting...")
	spot.SetQ(1e-3)
	if spot.Config().Q != 1e-3 {
		t.Errorf("Error while setting Q (expected 1e-3, got %f)", spot.Config().Q)
		testERROR()
	} else {
		testOK()
	}

	// data
	var N = 80000
	data := gaussianSample(N)

	for i := 0; i < N; i++ {
		spot.Step(data[i])
	}

	//
	var zTrue = 3.09
	relativeError := 100. * math.Abs(zTrue-spot.GetUpperThreshold()) / zTrue

	checkTitle("Checking error...")

	if relativeError > 5.0 {
		t.Error("Error")
		testERROR()
	} else if relativeError > 2.5 {
		testWARNING()
	}
	testOK()
	spot.Delete()
}

func TestSpotStatus(t *testing.T) {
	title("Testing Spot status")
	// init spot object
	var q = 1e-3
	var nInit int32 = 10000
	var level = 0.995
	up, down, alert, bounded := true, true, false, true
	var maxExcess int32 = 300

	spot := NewSpot(q, nInit, level, up, down, alert, bounded, maxExcess)

	var N = 15000
	data := gaussianSample(N)

	for i := 0; i < N; i++ {
		spot.Step(data[i])
	}

	status := spot.Status()
	checkTitle("Checking n...")
	if status.N != int32(N) {
		t.Error("Error on the number of observations")
		testERROR()
	} else {
		testOK()
	}

	checkTitle("Checking NtUp/NtDown...")
	eUp := status.NtUp - int32((1.-level)*float64(N))
	eDown := status.NtDown - int32((1.-level)*float64(N))
	if math.Abs(float64(eUp)) > 30 || math.Abs(float64(eDown)) > 30 {
		t.Error("Error on the number of excess")
		testERROR()
	} else {
		testOK()
	}

	checkTitle("Checking TUp/TDown...")
	if math.Abs(status.TUp-2.576) > 0.1 || math.Abs(status.TDown+2.576) > 0.1 {
		t.Error("Error on the TUp/TDown values")
		testERROR()
	} else {
		testOK()
	}

	checkTitle("Checking alarms...")
	if (status.AlDown != 0) || (status.AlUp != 0) {
		t.Error("Alarms have been triggered while alarm mode is desactivate")
		testERROR()
	} else {
		testOK()
	}

	fmt.Println("\n", status)
}

func TestSpotConfig(t *testing.T) {
	title("Testing Spot config")
	// init spot object
	var q = 1e-3
	var nInit int32 = 3000
	var level = 0.995
	up, down, alert, bounded := true, true, false, true
	var maxExcess int32 = 300

	spot := NewSpot(q, nInit, level, up, down, alert, bounded, maxExcess)

	config := spot.Config()

	checkTitle("Checking Q...")
	if config.Q != q {
		t.Error("Error about the value of q")
		testERROR()
	} else {
		testOK()
	}

	checkTitle("Checking Ninit...")
	if config.Ninit != nInit {
		t.Error("Error about the value of nInit")
	} else {
		testOK()
	}

	checkTitle("Checking Level...")
	if config.Level != level {
		t.Error("Error about the value of level")
		testERROR()
	} else {
		testOK()
	}

	checkTitle("Checking Up...")
	if config.Up != up {
		t.Error("Error about the value of up")
		testERROR()
	} else {
		testOK()
	}

	checkTitle("Checking Down...")
	if config.Down != down {
		t.Error("Error about the value of down")
		testERROR()
	} else {
		testOK()
	}

	checkTitle("Checking Alert...")
	if config.Alert != alert {
		t.Error("Error about the value of alert")
		testERROR()
	} else {
		testOK()
	}

	checkTitle("Checking Bounded...")
	if config.Bounded != bounded {
		t.Error("Error about the value of bounded")
		testERROR()
	} else {
		testOK()
	}

	checkTitle("Checking MaxExcess...")
	if config.MaxExcess != maxExcess {
		t.Error("Error about the value of maxExcess")
		testERROR()
	} else {
		testOK()
	}
	fmt.Println("\n", spot.Config())
}

func TestBasicSpotAccess(t *testing.T) {
	title("Testing Spot status accesses")

	config := SpotConfig{
		Q:         1e-4,
		Ninit:     5000,
		Level:     0.99,
		Up:        true,
		Down:      true,
		Alert:     false,
		Bounded:   true,
		MaxExcess: 200}

	spot := NewSpotFromConfig(config)

	N := int(config.Ninit)

	data := gaussianSample(N)

	for i := 0; i < N; i++ {
		spot.Step(data[i])
	}

	checkTitle("Checking TUp/TDown...")
	errTUp := math.Abs(spot.GetUpperT()-2.326) / 2.326
	errTDown := math.Abs(spot.GetLowerT()+2.326) / 2.326
	if errTUp > 5 || errTDown > 5 {
		t.Error("Error on the TUp/TDown values")
		fmt.Println(spot.Status())
		testERROR()
	} else if errTUp > 2.5 || errTDown > 2.5 {
		testWARNING()
	} else {
		testOK()
	}

	checkTitle("Checking ZUp/ZDown...")
	errZUp := math.Abs(spot.GetUpperThreshold()-3.719) / 3.719
	errZDown := math.Abs(spot.GetLowerThreshold()+3.719) / 3.719
	if errZUp > 5 || errZDown > 5 {
		t.Error("Error on the TUp/TDown values")
		fmt.Println("\n", spot.Status())
		testERROR()
	} else if errZUp > 2.5 || errZDown > 2.5 {
		testWARNING()
	} else {
		testOK()
	}
}

func TestSpotProbabilityComputation(t *testing.T) {
	title("Probability computation")

	config := SpotConfig{
		Q:         1e-4,
		Ninit:     10000,
		Level:     0.999,
		Up:        true,
		Down:      true,
		Alert:     false,
		Bounded:   true,
		MaxExcess: 200}

	spot := NewSpotFromConfig(config)
	N := int(config.Ninit)
	data := gaussianSample(N)

	for i := 0; i < N; i++ {
		spot.Step(data[i])
	}

	checkTitle("Checking Up probability computation...")
	errUp := math.Abs(spot.UpProbability(3.09)-1e-3) / 1e-3
	if errUp > 5 {
		testERROR()
		t.Errorf("Expected 1e-3, got %f", spot.UpProbability(3.09))
	} else if errUp > 2.5 {
		testWARNING()
	} else {
		testOK()
	}

	checkTitle("Checking Down probability computation...")
	errDown := math.Abs(spot.DownProbability(-3.09)-1e-3) / 1e-3
	if errDown > 5 {
		testERROR()
		t.Errorf("Expected 1e-3, got %f", spot.DownProbability(-3.09))
	} else if errDown > 2.5 {
		testWARNING()
	} else {
		testOK()
	}

	config = SpotConfig{
		Q:         1e-4,
		Ninit:     10000,
		Level:     0.999,
		Up:        false,
		Down:      false,
		Alert:     false,
		Bounded:   true,
		MaxExcess: 200}
	spot = NewSpotFromConfig(config)
	checkTitle("Checking NaN (Up)...")
	if !math.IsNaN(spot.UpProbability(12.)) {
		testERROR()
	} else {
		testOK()
	}

	checkTitle("Checking NaN (Down)...")
	if !math.IsNaN(spot.DownProbability(-12.)) {
		testERROR()
	} else {
		testOK()
	}
}
