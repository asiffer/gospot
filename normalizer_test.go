// normalizer_test.go
package gospot

import (
	"math/rand"
	"testing"
)

func TestNormalizerFeed(t *testing.T) {
	title("Testing feeding normalizer")
	// var z float64
	var err error
	depth := 10
	normalizer := NewNormalizer(depth, true, false)
	checkTitle("Checking size...")
	if normalizer.Depth() != depth {
		testERROR()
		t.Errorf("Expected %d, got %d", depth, normalizer.Depth())
	} else {
		testOK()
	}

	checkTitle("Checking transitory steps...")
	for i := 0; i < depth; i++ {
		_, err = normalizer.Step(rand.Float64())
		if err == nil {
			t.Error("Error in transitory steps")
			testERROR()
			return
		}
	}
	testOK()

	checkTitle("Checking cruising steps...")
	for i := 0; i < depth; i++ {
		_, err = normalizer.Step(rand.Float64())
		if err != nil {
			t.Error("Error in cruising steps")
			testERROR()
			return
		}
	}
	testOK()

}

func TestCentering(t *testing.T) {
	title("Testing centering")
	var z float64
	// var err error
	depth := 10
	val := 2.5
	normalizer := NewNormalizer(depth, true, false)
	for i := 0; i < depth; i++ {
		normalizer.Step(val)
	}

	checkTitle("Checking centered value...")
	for i := 0; i < 2*depth; i++ {
		z, _ = normalizer.Step(val)
		if z != 0.0 {
			t.Error("Error while centering")
			testERROR()
			return
		}
	}
	testOK()
}

func TestScaling(t *testing.T) {
	title("Testing scaling")
	var z float64
	depth := 20

	normalizer := NewNormalizer(2*depth, true, false)
	for i := 0; i < depth; i++ {
		normalizer.Step(0.0)
	}
	for i := 0; i < depth; i++ {
		normalizer.Step(2.0)
	}

	z, _ = normalizer.Step(17.0)
	checkTitle("Checking scaled value...")
	if z != 16.0 {
		t.Error("Error while scaling")
		testERROR()
	} else {
		testOK()
	}
}
