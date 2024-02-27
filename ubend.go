package gospot

import (
	"math"
)

type Ubend struct {
	Cursor         uint64    `json:"cursor"`
	Capacity       uint64    `json:"capacity"`
	LastErasedData float64   `json:"last_erased_data"`
	Filled         bool      `json:"filled"`
	Data           []float64 `json:"data"`
}

func NewUbend(capacity uint64) *Ubend {
	return &Ubend{
		Cursor:         0,
		Capacity:       capacity,
		LastErasedData: math.NaN(),
		Filled:         false,
		Data:           make([]float64, capacity),
	}
}

func (ubend *Ubend) Size() uint64 {
	if ubend.Filled {
		return ubend.Capacity
	}
	return ubend.Cursor
}

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
