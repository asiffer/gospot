package gospot

import (
	"math"
)

type Ubend struct {
	// Current position inside the container
	Cursor uint64 `json:"cursor"`
	// Max storage
	Capacity uint64 `json:"capacity"`
	// Last erased value (i.e. replaced by a new one)
	LastErasedData float64 `json:"last_erased_data"`
	// Container fill status
	Filled bool `json:"filled"`
	// Data container
	Data []float64 `json:"data"`
}

// NewUbend initializes a new [Ubend] structure given a max capacity
func NewUbend(capacity uint64) *Ubend {
	return &Ubend{
		Cursor:         0,
		Capacity:       capacity,
		LastErasedData: math.NaN(),
		Filled:         false,
		Data:           make([]float64, capacity),
	}
}

// Size returns the current size of the container
func (ubend *Ubend) Size() uint64 {
	if ubend.Filled {
		return ubend.Capacity
	}
	return ubend.Cursor
}

// Push a new value to the container and returns the erased one (or NaN if it does not exist)
func (ubend *Ubend) Push(x float64) float64 {
	if ubend.Filled {
		ubend.LastErasedData = ubend.Data[ubend.Cursor]
	}

	ubend.Data[ubend.Cursor] = x

	if ubend.Cursor == (ubend.Capacity - 1) {
		ubend.Cursor = 0
		ubend.Filled = true
	} else {
		ubend.Cursor++
	}

	return ubend.LastErasedData
}
