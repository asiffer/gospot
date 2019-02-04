// normalizer_test.go
package gospot

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestNormalizer(t *testing.T) {
	title("NORMALIZER")
}
func TestNormalizerFeed(t *testing.T) {
	fmt.Println("** Testing feeding normalizer **")
	// var z float64
	var err error
	depth := 10
	normalizer := NewNormalizer(depth, true, false)
	for i := 0; i < depth; i++ {
		_, err = normalizer.Step(rand.Float64())
		if err == nil {
			t.Error("Error in transitory steps")
		}
	}

	for i := 0; i < depth; i++ {
		_, err = normalizer.Step(rand.Float64())
		if err != nil {
			t.Error("Error in cruising steps")
		}
	}
}

func TestCentering(t *testing.T) {
	fmt.Println("** Testing centering **")
	var z float64
	// var err error
	depth := 10
	val := 2.5
	normalizer := NewNormalizer(depth, true, false)
	for i := 0; i < depth; i++ {
		normalizer.Step(val)
	}

	for i := 0; i < 2*depth; i++ {
		z, _ = normalizer.Step(val)
		if z != 0.0 {
			t.Error("Error while centering")
		}
	}
}

func TestScaling(t *testing.T) {
	fmt.Println("** Testing scaling **")
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
	if z != 16.0 {
		t.Error("Error while scaling")
	}
}
