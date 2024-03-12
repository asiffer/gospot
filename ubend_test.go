package gospot

import (
	"math"
	"math/rand"
	"testing"
)

func TestSize(t *testing.T) {
	var size uint64 = 20
	u := NewUbend(size)

	for i := uint64(0); i < size; i++ {
		u.Push(rand.Float64())
		if u.Size() != i+1 {
			t.Errorf("bad size: %d instead of %d", u.Size(), i+1)
		}
	}

	for i := uint64(0); i < size; i++ {
		u.Push(rand.Float64())
		if u.Size() != size {
			t.Errorf("bad size: %d instead of %d", u.Size(), size)
		}
	}
}

func TestCursor(t *testing.T) {
	var size uint64 = 20
	u := NewUbend(size)
	for i := uint64(0); i < size-1; i++ {
		u.Push(rand.Float64())
		c := (i + 1) % size
		if u.Cursor != c {
			t.Errorf("bad cursor: %v instead of %d", u.Cursor, c)
		}
		if u.Filled {
			t.Errorf("must not be filled")
		}
	}

	u.Push(rand.Float64())
	if u.Cursor != 0 {
		t.Errorf("cursor has not been reset")
	}
	if !u.Filled {
		t.Errorf("container must be filled")
	}
}

func TestLastErasedData(t *testing.T) {
	var size uint64 = 20
	u := NewUbend(size)
	if !math.IsNaN(u.LastErasedData) {
		t.Errorf("LastErasedData must be NaN")
	}

	for i := uint64(0); i < size; i++ {
		x := u.Push(float64(i))
		if !math.IsNaN(x) {
			t.Errorf("LastErasedData must be NaN")
		}
	}

	for i := uint64(0); i < size; i++ {
		led := u.Push(float64(2 * i))
		if led != float64(i) {
			t.Errorf("bad LastErasedData: %v != %v", led, i)
		}
	}
}
