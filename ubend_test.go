// ubend_test.go
package gospot

import (
	"fmt"
	"math"
	"testing"
)

func TestUbend(t *testing.T) {
	title("UBEND")
}

func TestCreateUbend(t *testing.T) {
	fmt.Println("** Testing Ubend initialization **")
	size := 10
	ubend := NewUbend(size)
	if ubend.Size() != size {
		t.Error("Bad size")
	}
	if ubend.Length() != 0 {
		t.Error("Bad length")
	}
	if !math.IsNaN(ubend.Mean()) {
		t.Error("Bad mean computation")
	}

}

func TestUbendPush(t *testing.T) {
	fmt.Println("** Testing Ubend push **")
	size := 10
	ubend := NewUbend(size)

	for i := 0; i < 8; i++ {
		ubend.Push(float64(i))
	}

	if ubend.Length() != 8 {
		t.Error("Bad length before reaching size")
	}

	for i := 8; i < 15; i++ {
		ubend.Push(float64(i))
	}

	if ubend.Length() != size {
		t.Error("Bad length after reaching size")
	}

	for i := 15; i < 42; i++ {
		ubend.Push(float64(i))
	}

}

func TestUbendMomentComputation(t *testing.T) {
	fmt.Println("** Testing Ubend moment computation **")
	size := 10
	val := 1.0
	ubend := NewUbend(size)

	for i := 0; i < 10; i++ {
		ubend.Push(val)
	}

	if ubend.Mean() != val {
		t.Error("Bad mean computation when container is filled")
	}

	if ubend.Var() != 0.0 {
		t.Error("Bad variance computation when container is filled")
	}

	for i := 0; i < 10; i++ {
		ubend.Push(2. * val)
	}

	if ubend.Mean() != 2.*val {
		t.Error("Bad mean computation in cruising regime")
	}

}
