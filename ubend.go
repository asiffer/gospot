// ubend.go

package gospot

import (
	"math"
)

// Ubend is a circular container of a given size. It appends
// new data until the size is reached. After that it replace
// older data with the new incoming ones.
type Ubend struct {
	data           []float64
	m              float64
	m2             float64
	id             int
	length         int
	size           int
	lastErasedData float64
}

// NewUbend creates a new Ubend structure.
func NewUbend(size int) *Ubend {
	return &Ubend{
		data:           make([]float64, 0),
		m:              0.0,
		m2:             0.0,
		id:             0,
		length:         0,
		size:           size,
		lastErasedData: math.NaN()}
}

// Length returns the current number of data in the container
func (u *Ubend) Length() int {
	// fmt.Println(len(u.data))
	return u.length
}

// Size return the capacity of the container, that is to say
// the maximum of data it can stores.
func (u *Ubend) Size() int {
	return u.size
}

// Clear resets the container. It keeps the original size.
func (u *Ubend) Clear() {
	u.data = make([]float64, 0)
	u.m = 0.0
	u.m2 = 0.0
	u.id = 0
	u.length = 0
}

// Push add a new data to the container. It updates the
// basic moments.
func (u *Ubend) Push(x float64) {
	if u.Length() < u.Size() || u.Size() <= 0 {
		u.data = append(u.data, x)
		u.length++
		// update moment
		u.m += x
		u.m2 += x * x
	} else if u.Size() > 0 {
		old := u.data[u.id]
		u.m -= old
		u.m2 -= old * old
		u.data[u.id] = x
		u.id = (u.id + 1) % u.size
		u.lastErasedData = old
		// update moment
		u.m += x
		u.m2 += x * x
	}
	// otherwise it means that Size == 0
	// so nothing is done...
}

// Cancel goes to the state before the last push
// Warning: you can cancel only one push. If you try more
// the container will be corrupted.
func (u *Ubend) Cancel() {
	if u.Length() > 0 {
		if math.IsNaN(u.lastErasedData) {
			old := u.data[u.Length()-1]
			u.data = u.data[:u.Length()-1]
			u.length--
			u.m -= old
			u.m2 -= old * old
		} else {
			// data to re-add
			old := u.lastErasedData
			// backstep
			u.id = (u.size + u.id - 1) % u.size
			// data to remove
			remove := u.data[u.id]
			// fmt.Printf("Id: %d, Remove: %f, Old: %f\n", u.id, remove, old)
			// Update
			u.m = u.m - remove + old
			// u.m -= old
			u.m2 = u.m2 - remove*remove + old*old
			// step forward
			u.id = (u.id + 1) % u.size
			// u.Push(old)
		}
	}
}

// IsFull returns whether the container is full (cruising regime)
// or not (transitory regime).
func (u *Ubend) IsFull() bool {
	return u.Size() == u.Length() || u.Size() <= 0
}

// Mean computes the mean of the current data
// of the container
func (u *Ubend) Mean() float64 {
	if u.Size() == 0 {
		return 0.0
	}
	return u.m / float64(u.Length())
}

// MeanSquare computes the mean of the square of
// the current data of the container
func (u *Ubend) MeanSquare() float64 {
	return u.m2 / float64(u.Length())
}

// Var computes the variance of the current data
// of the container
func (u *Ubend) Var() float64 {
	mean := u.Mean()
	return u.MeanSquare() - mean*mean
}

// Std computes the standard deviation of the current data
// of the container
func (u *Ubend) Std() float64 {
	return math.Sqrt(u.Var())
}
