// status.go

package gospot

import (
	"fmt"
)

// SpotStatus is the structure embedding the status of a Spot instance
type SpotStatus struct {
	// N is the number of normal observations (not the alarms)
	N int32
	// ExUp is the urrent number of up excesses
	ExUp int32
	// ExDown is the current number of down excesses
	ExDown int32
	// NtUp is the total number of up excesses
	NtUp int32
	// NtDown is the total number of down excesses
	NtDown int32
	// AlUp is the number of up alarms
	AlUp int32
	// AlDown is the number of down alarms
	AlDown int32
	// TUp is the transitional up threshold
	TUp float64
	// TDown is the transitional down threshold
	TDown float64
	// ZUp is the up alert thresholds
	ZUp float64
	// ZDown is the down alert thresholds
	ZDown float64
}

// DSpotStatus is the structure embedding the status of a DSpot instance
type DSpotStatus struct {
	SpotStatus
	// Mean is the the value of the current local model
	Mean float64
}

var (
	spotStatusNew    func(uintptr) uintptr
	spotStatusDelete func(uintptr)
	statusGetN       func(uintptr) int32
	statusGetExUp    func(uintptr) int32
	statusGetExDown  func(uintptr) int32
	statusGetNtUp    func(uintptr) int32
	statusGetNtDown  func(uintptr) int32
	statusGetAlUp    func(uintptr) int32
	statusGetAlDown  func(uintptr) int32
	statusGetTUp     func(uintptr) float64
	statusGetTDown   func(uintptr) float64
	statusGetZUp     func(uintptr) float64
	statusGetZDown   func(uintptr) float64
)

var statusSymbols = map[string]interface{}{
	"Spot_status_ptr":     &spotStatusNew,
	"Spot_status_delete":  &spotStatusDelete,
	"_status_get_n":       &statusGetN,
	"_status_get_ex_up":   &statusGetExUp,
	"_status_get_ex_down": &statusGetExDown,
	"_status_get_Nt_up":   &statusGetNtUp,
	"_status_get_Nt_down": &statusGetNtDown,
	"_status_get_al_up":   &statusGetAlUp,
	"_status_get_al_down": &statusGetAlDown,
	"_status_get_t_up":    &statusGetTUp,
	"_status_get_t_down":  &statusGetTDown,
	"_status_get_z_up":    &statusGetZUp,
	"_status_get_z_down":  &statusGetZDown,
}

// LoadSymbolsStatus It loads the symbols related to the SpotStatus object
// from the C++ library libspot. It returns an error if a loading fails, nil
// pointer otherwise
func LoadSymbolsStatus() error {
	var err error
	for k, f := range statusSymbols {
		err = libspot.Sym(k, f)
		if err != nil {
			return fmt.Errorf("Error in loading %s (%s)", k, err.Error())
		}
	}
	return nil
}

func (ss SpotStatus) String() string {
	return fmt.Sprintf("%8s %d\n%8s %d\n%8s %d\n%8s %d\n%8s %d\n%8s %d\n%8s %d\n%8s %.6f\n%8s %.6f\n%8s %.6f\n%8s %.6f\n",
		"n", ss.N, "ex_up", ss.ExUp, "ex_down", ss.ExDown, "Nt_up", ss.NtUp, "Nt_down", ss.NtDown, "al_up", ss.AlUp, "al_down", ss.AlDown, "t_up", ss.TUp, "t_down", ss.TDown, "z_up", ss.ZUp, "z_down", ss.ZDown)
}

func (dss DSpotStatus) String() string {
	return fmt.Sprintf("%8s %.6f\n%8s %d\n%8s %d\n%8s %d\n%8s %d\n%8s %d\n%8s %d\n%8s %d\n%8s %.6f\n%8s %.6f\n%8s %.6f\n%8s %.6f\n",
		"drift", dss.Mean,
		"n", dss.N,
		"ex_up", dss.ExUp,
		"ex_down", dss.ExDown,
		"Nt_up", dss.NtUp,
		"Nt_down", dss.NtDown,
		"al_up", dss.AlUp,
		"al_down", dss.AlDown,
		"t_up", dss.TUp,
		"t_down", dss.TDown,
		"z_up", dss.ZUp,
		"z_down", dss.ZDown)
}
