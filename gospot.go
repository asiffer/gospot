// gospot.go
package gospot

import (
	"C"

	"github.com/rainycape/dl"
)
import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var libspot *dl.DL
var LOCAL_LIB string = "/home/asr/Documents/Work/github/libspot/master/lib/libspot.so"
var IsInitialized bool
var logger zerolog.Logger // package logger

func bool2uint8(b bool) uint8 {
	if b {
		return uint8(1)
	} else {
		return uint8(0)
	}
}

type SpotInterface interface {
	Delete()
	Step(x float64) int32
	GetUpperT() float64
	GetLowerT() float64
	GetUpperThreshold() float64
	GetLowerThreshold() float64
	SetQ(q float64)
	UpProbability(z float64) float64
	DownProbability(z float64) float64
}

func init() {
	var err error
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.StampMicro}
	logger = zerolog.New(output).With().Timestamp().Logger()
	// lib, err := dl.Open("libspot", dl.RTLD_LOCAL)
	lib, err := dl.Open(LOCAL_LIB, dl.RTLD_LOCAL)
	if err != nil {
		logger.Error().Msgf("Error while loading libspot (%s).\nMaybe the library is not installed on your system, see https://asiffer.github.io/libspot/", err)
		return
	}
	libspot = lib
	// fmt.Println("Library found")

	if err = LoadSymbolsSpot(); err != nil {
		logger.Error().Msgf("Error while loading Spot symbols (%s)", err)
		return
	}

	if err = LoadSymbolsStatus(); err != nil {
		logger.Error().Msgf("Error while loading Status symbols (%s)", err)
		return
	}

	if err = LoadSymbolsConfig(); err != nil {
		logger.Error().Msgf("Error while loading Config symbols (%s)", err)
		return
	}

	DefaultSpotConfig = SpotConfig{Q: 1e-4, N_init: 1000, Level: 0.99, Up: true, Down: true, Alert: true, Bounded: true, Max_excess: 200}
	DefaultDSpotConfig = DSpotConfig{Depth: 0, SpotConfig: DefaultSpotConfig}
	IsInitialized = true
}

func main() {
}
