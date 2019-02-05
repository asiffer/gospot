// ubend_test.go
package gospot

import (
	"math"
	"testing"
)

func TestCreateUbend(t *testing.T) {
	title("Testing Ubend initialization")
	size := 10
	ubend := NewUbend(size)

	checkTitle("Checking size...")
	if ubend.Size() != size {
		t.Error("Bad size")
		testERROR()
	} else {
		testOK()
	}

	checkTitle("Checking length...")
	if ubend.Length() != 0 {
		t.Error("Bad length")
		testERROR()
	} else {
		testOK()
	}

	checkTitle("Checking mean computation...")
	if !math.IsNaN(ubend.Mean()) {
		t.Error("Bad mean computation")
		testERROR()
	} else {
		testOK()
	}

}

func TestUbendPush(t *testing.T) {
	title("Testing Ubend push")
	size := 10
	ubend := NewUbend(size)

	for i := 0; i < 8; i++ {
		ubend.Push(float64(i))
	}

	checkTitle("Checking length (<size)...")
	if ubend.Length() != 8 {
		t.Error("Bad length before reaching size")
		testERROR()
	} else {
		testOK()
	}

	for i := 8; i < 15; i++ {
		ubend.Push(float64(i))
	}

	checkTitle("Checking length (>size)...")
	if ubend.Length() != size {
		t.Error("Bad length after reaching size")
		testERROR()
	} else {
		testOK()
	}

	for i := 15; i < 42; i++ {
		ubend.Push(float64(i))
	}

}

func TestUbendMomentComputation(t *testing.T) {
	title("Testing Ubend moment computation")
	size := 10
	val := 1.0
	ubend := NewUbend(size)

	for i := 0; i < 10; i++ {
		ubend.Push(val)
	}

	checkTitle("Checking mean computation (filled)...")
	if ubend.Mean() != val {
		t.Error("Bad mean computation when container is filled")
		testERROR()
	} else {
		testOK()
	}

	checkTitle("Checking variance computation (filled)...")
	if ubend.Var() != 0.0 {
		t.Error("Bad variance computation when container is filled")
		testERROR()
	} else {
		testOK()
	}
	for i := 0; i < 10; i++ {
		ubend.Push(2. * val)
	}

	checkTitle("Checking mean computation (cruising regime)...")
	if ubend.Mean() != 2.*val {
		t.Error("Bad mean computation in cruising regime")
		testERROR()
	} else {
		testOK()
	}

}
