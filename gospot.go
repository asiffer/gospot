// gospot.go

// Package gospot re-implements libspot
package gospot

var (
	// DefaultSpotConfig is the default structure to create a Spot Object
	DefaultSpotConfig = SpotConfig{
		Q:         1e-4,
		Ninit:     1500,
		Level:     0.98,
		Up:        true,
		Down:      true,
		Alert:     true,
		Bounded:   true,
		MaxExcess: 200}
	// DefaultDSpotConfig is the default structure to create a DSpot Object
	DefaultDSpotConfig = DSpotConfig{
		Depth:      0,
		SpotConfig: DefaultSpotConfig}
)

func main() {

}
