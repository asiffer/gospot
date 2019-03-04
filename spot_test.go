// spot_test.go

package gospot

import (
	"fmt"
	"math"
	"testing"
)

func TestInitSpot(t *testing.T) {
	title("Spot")
	// rand.Seed(time.Now().UnixNano())
}
func TestSpotBasicRun(t *testing.T) {
	checkTitle("Basic run...")
	spot := NewDefaultSpot()
	data := standardGaussianSample(8000)
	for _, x := range data {
		spot.Step(x)
	}
	testOK()
}

func TestSpotThresholdComputation(t *testing.T) {
	title("Testing Spot threshold computation")

	sc := SpotConfig{
		Q:         1e-4,
		Ninit:     4000,
		Level:     0.99,
		Up:        true,
		Down:      true,
		Alert:     false,
		Bounded:   true,
		MaxExcess: 200}

	spot := NewSpotFromConfig(&sc)

	checkTitle("Checking Q setting...")
	spot.SetQ(1e-3)
	if spot.Config().Q != 1e-3 {
		t.Errorf("Error while setting Q (expected 1e-3, got %f)", spot.Config().Q)
		testERROR()
	} else {
		testOK()
	}

	// data
	var N = 5000
	data := standardGaussianSample(N)

	for i := 0; i < N; i++ {
		spot.Step(data[i])
	}

	var zTrue = 3.09
	relativeError := 100. * math.Abs(zTrue-spot.GetUpperThreshold()) / zTrue

	checkTitle("Checking error...")

	if relativeError > 5.0 {
		t.Errorf("Expected lower than 5%%, got %.2f", relativeError)
		testERROR()
		fmt.Println(spot.Status())
	} else if relativeError > 2.5 {
		testWARNING()
	} else {
		testOK()
	}

}

func TestSpotProbabilityComputation(t *testing.T) {
	title("Spot Probability computation")

	config := SpotConfig{
		Q:         1e-4,
		Ninit:     10000,
		Level:     0.999,
		Up:        true,
		Down:      true,
		Alert:     false,
		Bounded:   true,
		MaxExcess: 200}

	spot := NewSpotFromConfig(&config)
	N := config.Ninit
	data := standardGaussianSample(N)

	for i := 0; i < N; i++ {
		spot.Step(data[i])
	}

	checkTitle("Checking Up probability computation...")
	errUp := math.Abs(spot.UpProbability(3.09)-1e-3) / 1e-3
	if errUp > 5 {
		testERROR()
		t.Errorf("Expected 1e-3, got %f", spot.UpProbability(3.09))
		fmt.Println(spot.up)
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
		fmt.Println(spot.down)
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
	spot = NewSpotFromConfig(&config)
	checkTitle("Checking NaN (Up)...")
	if math.IsNaN(spot.UpProbability(12.)) && math.IsNaN(spot.GetUpperT()) && math.IsNaN(spot.GetUpperThreshold()) {
		testOK()
	} else {
		testERROR()
	}

	checkTitle("Checking NaN (Down)...")
	if math.IsNaN(spot.DownProbability(12.)) && math.IsNaN(spot.GetLowerT()) && math.IsNaN(spot.GetLowerThreshold()) {
		testOK()
	} else {
		testERROR()
	}
}
