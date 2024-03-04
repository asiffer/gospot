# gospot

![Test](https://github.com/asiffer/gospot/workflows/Test/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/asiffer/gospot)](https://goreportcard.com/report/github.com/asiffer/gospot)
[![Coverage Status](https://codecov.io/github/asiffer/gospot/coverage.svg?branch=master)](https://codecov.io/github/asiffer/gospot?branch=master)
[![GoDoc](https://godoc.org/github.com/asiffer/gospot?status.svg)](https://godoc.org/github.com/asiffer/gospot)

`gospot` is a pure golang implementation of [libspot](https://asiffer.github.io/libspot/).
This module roughly follows the same structure.

> [!CAUTION]
> The last version (v0.2) includes many breaking changes. If your project cannot be migrated, you can still points to
> the previous one: `go get github.com/asiffer/gospot@v0.1.1`

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

func gaussian(size uint64) <-chan float64 {
	out := make(chan float64, 1)
	go func() {
		for i := uint64(0); i < size; i++ {
			out <- rand.NormFloat64()
		}
		close(out)
	}()
	return out
}

func gaussianBatch(size uint64) []float64 {
	out := make([]float64, size)
	k := 0
	for x := range gaussian(size) {
		out[k] = x
		k++
	}
	return out
}

func main() {
	s, err := gospot.NewSpot(1e-5, false, true, 0.99, 2000)
	if err != nil {
		panic(err)
	}
	data := gaussianBatch(10000)
	s.Fit(data)

	A := 0
	E := 0
	N := 0

	for x := range gaussian(1000000) {
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
