# gospot

![Test](https://github.com/asiffer/gospot/workflows/Test/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/asiffer/gospot)](https://goreportcard.com/report/github.com/asiffer/gospot)
[![Coverage Status](https://codecov.io/github/asiffer/gospot/coverage.svg?branch=master)](https://codecov.io/github/asiffer/gospot?branch=master)
[![GoDoc](https://godoc.org/github.com/asiffer/gospot?status.svg)](https://godoc.org/github.com/asiffer/gospot)

`gospot` is a pure golang implementation of [libspot](https://asiffer.github.io/libspot/).
This module roughly follows the same structure.

## Download

```shell
$ go get github.com/asiffer/gospot
```

## Usage

Once `gospot` is imported, you can create a `Spot` object and feed some data.

```golang
// example/main.go
package main

import (
	"math/rand"

	"github.com/asiffer/gospot"
)

func gaussian(size uint64) []float64 {
	out := make([]float64, size)
	for i := uint64(0); i < size; i++ {
		out[i] = rand.NormFloat64()
	}
	return out
}

func main() {
	s, _ := gospot.NewSpot(1e-5, false, true, 0.99, 2000)
	data := gaussian(10000)
	s.Fit(data)

	A := 0
	E := 0
	N := 0

	for _, x := range gaussian(1000000) {
		switch s.Step(x) {
		case gospot.ANOMALY:
			A++
		case gospot.EXCESS:
			E++
		default:
			N++
		}
	}
}

```
