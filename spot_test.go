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

func TestSpot(t *testing.T) {
	title("SPOT")
}

func TestInitAndRun(t *testing.T) {
	fmt.Println("** Testing initialization and run **")
	// init spot object
	var q float64 = 1e-4
	var n_init int32 = 2000
	var level float64 = 0.99
	up, down, alert, bounded := true, true, true, true
	var max_excess int32 = 200

	spot := NewSpot(q, n_init, level, up, down, alert, bounded, max_excess)
	// data
	var N int = 10000
	data := gaussianSample(N)

	for i := 0; i < N; i++ {
		spot.Step(data[i])
	}

	spot.Delete()
}

func TestThresholdComputation(t *testing.T) {
	fmt.Println("** Testing threshold computation **")
	// init spot object
	var q float64 = 1e-3
	var n_init int32 = 3000
	var level float64 = 0.99
	up, down, alert, bounded := true, true, false, true
	var max_excess int32 = 300

	spot := NewSpot(q, n_init, level, up, down, alert, bounded, max_excess)

	// data
	var N int = 80000
	data := gaussianSample(N)

	for i := 0; i < N; i++ {
		spot.Step(data[i])
	}

	//
	var z_true float64 = 3.09
	relativeError := 100. * math.Abs(z_true-spot.GetUpperThreshold()) / z_true
	fmt.Printf("Relative error: %.2f%%\n", relativeError)
	if relativeError > 5.0 {
		t.Error("Warning")
	}
	spot.Delete()
}

func TestStatus(t *testing.T) {
	fmt.Println("** Testing status **")
	// init spot object
	var q float64 = 1e-3
	var n_init int32 = 3000
	var level float64 = 0.995
	up, down, alert, bounded := true, true, false, true
	var max_excess int32 = 300

	spot := NewSpot(q, n_init, level, up, down, alert, bounded, max_excess)

	var N int = 15000
	data := gaussianSample(N)

	for i := 0; i < N; i++ {
		spot.Step(data[i])
	}

	status := spot.Status()
	fmt.Println(status)
	if (status.Al_down != 0) || (status.Al_up != 0) {
		t.Error("Alarms have been triggered while alarm mode is desactivate")
	}
}

func TestConfig(t *testing.T) {
	fmt.Println("** Testing config **")
	// init spot object
	var q float64 = 1e-3
	var n_init int32 = 3000
	var level float64 = 0.995
	up, down, alert, bounded := true, true, false, true
	var max_excess int32 = 300

	spot := NewSpot(q, n_init, level, up, down, alert, bounded, max_excess)

	config := spot.Config()

	if config.Q != q {
		t.Error("Error about the value of q")
	}
	if config.N_init != n_init {
		t.Error("Error about the value of n_init")
	}
	if config.Level != level {
		t.Error("Error about the value of level")
	}
	if config.Up != up {
		t.Error("Error about the value of up")
	}
	if config.Down != down {
		t.Error("Error about the value of down")
	}
	if config.Alert != alert {
		t.Error("Error about the value of alert")
	}
	if config.Bounded != bounded {
		t.Error("Error about the value of bounded")
	}
	if config.Max_excess != max_excess {
		t.Error("Error about the value of max_excess")
	}
	fmt.Println(spot.Config())
}
