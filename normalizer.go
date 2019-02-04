// normalizer.go
package gospot

import (
	"errors"
	"math"
)

type Normalizer struct {
	ubend  *Ubend
	center bool
	scale  bool
}

// NewNormalizer creates a new Normalizer instance.
// depth is the size of the underlying moving average
// centering means that data will be centered
// scaling means that data will be divided by the standard deviation
func NewNormalizer(depth int, centering bool, scaling bool) *Normalizer {
	if depth == 0 {
		centering = false
		scaling = false
	}
	return &Normalizer{ubend: NewUbend(depth),
		center: centering,
		scale:  scaling}
}

// Depth returns the depth the underlying moving average
func (n *Normalizer) Depth() int {
	return n.ubend.Size()
}

// Average returns the value of the underlying moving average
func (n *Normalizer) Average() float64 {
	return n.ubend.Mean()
}

// Step returns a normalized version of the new incoming value x.
// It then stores x to update
func (n *Normalizer) Step(x float64) (float64, error) {
	if n.ubend.IsFull() {
		var z float64 = x
		if n.center {
			z = z - n.ubend.Mean()
		}
		if n.scale {
			z = z / n.ubend.Std()
		}
		n.ubend.Push(x)
		return z, nil
	} else {
		n.ubend.Push(x)
		return math.NaN(), errors.New("The depth is not reached yet")
	}
}

// Cancel removes the last step only
func (n *Normalizer) Cancel() {
	n.ubend.Cancel()
}
