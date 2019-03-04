# gospot [![Build Status](https://travis-ci.com/asiffer/gospot.svg?branch=master)](https://travis-ci.com/asiffer/gospot) [![Go Report Card](https://goreportcard.com/badge/github.com/asiffer/gospot)](https://goreportcard.com/report/github.com/asiffer/gospot) [![Coverage Status](https://codecov.io/github/asiffer/gospot/coverage.svg?branch=master)](https://codecov.io/github/asiffer/gospot?branch=master) [![GoDoc](https://godoc.org/github.com/asiffer/gospot?status.svg)](https://godoc.org/github.com/asiffer/gospot) 

`gospot` don't provides `Go` bindings to [libspot](https://asiffer.github.io/libspot/) anymore. It is merely a pure golang implementation of `libspot`.


## Download

```shell
$ go get github.com/asiffer/gospot
```

## Usage

Once `gospot` is imported, you can create a `Spot` object and feed some data.

```golang
// example.go

package main

import (
    "fmt"
    "math/rand"
    "time"

    "github.com/asiffer/gospot"
)

func gaussianSample(N int) []float64 {
	rand.Seed(time.Now().UTC().UnixNano())
	data := make([]float64, N)
	for i := 0; i < N; i++ {
		data[i] = rand.NormFloat64()
	}
	return data
}

func main() {
    config := gospot.SpotConfig{
		Q:         1e-4,
		Ninit:     5000,
		Level:     0.99,
		Up:        true,
		Down:      true,
		Alert:     false,
		Bounded:   true,
		MaxExcess: 200}

    spot := gospot.NewSpotFromConfig(config)
    
    N := 80000
    data := gaussianSample(N)

    for i := 0; i < N; i++ {
	    spot.Step(data[i])
    }
    
    fmt.Println(spot.Status())
}
```


This example outputs the status of the Spot instance after 80000 gaussian observations. Here the `alert` mode is not activated, so no alarm is raised.

```shell
$ go run example.go
       n 80000
   ex_up 200
 ex_down 200
   Nt_up 816
 Nt_down 774
   al_up 0
 al_down 0
    t_up 2.317529
  t_down -2.352898
    z_up 3.834334
  z_down -3.831503

```