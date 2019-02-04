// status.go
package gospot

import (
	"errors"
	"fmt"
)

type SpotStatus struct {
	// n Number of normal observations (not the alarms)
	N int32
	// ex_up Current number of up excesses
	Ex_up int32
	// ex_up Current number of down excesses
	Ex_down int32
	// Nt_up Total number of up excesses
	Nt_up int32
	// Total number of down excesses
	Nt_down int32
	// al_up Number of up alarms
	Al_up int32
	// al_down Number of down alarms
	Al_down int32
	// t_up transitional up threshold
	T_up float64
	// t_down transitional down threshold
	T_down float64
	// z_up up alert thresholds
	Z_up float64
	// z_down down alert thresholds
	Z_down float64
}

type DSpotStatus struct {
	SpotStatus
	// Mean the is the value of the current lcoal model
	Mean float64
}

var (
	spot_status_new    func(uintptr) uintptr
	spot_status_delete func(uintptr)
	status_get_n       func(uintptr) int32
	status_get_ex_up   func(uintptr) int32
	status_get_ex_down func(uintptr) int32
	status_get_Nt_up   func(uintptr) int32
	status_get_Nt_down func(uintptr) int32
	status_get_al_up   func(uintptr) int32
	status_get_al_down func(uintptr) int32
	status_get_t_up    func(uintptr) float64
	status_get_t_down  func(uintptr) float64
	status_get_z_up    func(uintptr) float64
	status_get_z_down  func(uintptr) float64
)

// LoadSymbolsStatus It loads the symbols related to the SpotStatus object
// from the C++ library libspot. It returns an error if a loading fails, nil
// pointer otherwise
func LoadSymbolsStatus() error {
	var err error
	err = libspot.Sym("Spot_status_ptr", &spot_status_new)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading Spot_status_ptr(%s)", err.Error()))
	}
	err = libspot.Sym("Spot_status_delete", &spot_status_delete)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading Spot_status_delete (%s)", err.Error()))
	}
	err = libspot.Sym("_status_get_n", &status_get_n)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading status_get_n (%s)", err.Error()))
	}
	err = libspot.Sym("_status_get_ex_up", &status_get_ex_up)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading status_get_ex_up (%s)", err.Error()))
	}
	err = libspot.Sym("_status_get_ex_down", &status_get_ex_down)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading status_get_ex_down (%s)", err.Error()))
	}
	err = libspot.Sym("_status_get_Nt_up", &status_get_Nt_up)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading status_get_Nt_up (%s)", err.Error()))
	}
	err = libspot.Sym("_status_get_Nt_down", &status_get_Nt_down)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading status_get_Nt_down (%s)", err.Error()))
	}
	err = libspot.Sym("_status_get_al_up", &status_get_al_up)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading status_get_al_up (%s)", err.Error()))
	}
	err = libspot.Sym("_status_get_al_down", &status_get_al_down)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading status_get_al_down (%s)", err.Error()))
	}
	err = libspot.Sym("_status_get_t_up", &status_get_t_up)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading status_get_t_up (%s)", err.Error()))
	}
	err = libspot.Sym("_status_get_t_down", &status_get_t_down)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading status_get_t_down (%s)", err.Error()))
	}
	err = libspot.Sym("_status_get_z_up", &status_get_z_up)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading status_get_z_up (%s)", err.Error()))
	}
	err = libspot.Sym("_status_get_z_down", &status_get_z_down)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in loading status_get_z_down (%s)", err.Error()))
	}
	return nil
}

func (ss SpotStatus) String() string {
	return fmt.Sprintf("%8s %d\n%8s %d\n%8s %d\n%8s %d\n%8s %d\n%8s %d\n%8s %d\n%8s %.6f\n%8s %.6f\n%8s %.6f\n%8s %.6f\n",
		"n", ss.N, "ex_up", ss.Ex_up, "ex_down", ss.Ex_down, "Nt_up", ss.Nt_up, "Nt_down", ss.Nt_down, "al_up", ss.Al_up, "al_down", ss.Al_down, "t_up", ss.T_up, "t_down", ss.T_down, "z_up", ss.Z_up, "z_down", ss.Z_down)
}

func (dss DSpotStatus) String() string {
	return fmt.Sprintf("%8s %.6f\n%8s %d\n%8s %d\n%8s %d\n%8s %d\n%8s %d\n%8s %d\n%8s %d\n%8s %.6f\n%8s %.6f\n%8s %.6f\n%8s %.6f\n",
		"drift", dss.Mean,
		"n", dss.N,
		"ex_up", dss.Ex_up,
		"ex_down", dss.Ex_down,
		"Nt_up", dss.Nt_up,
		"Nt_down", dss.Nt_down,
		"al_up", dss.Al_up,
		"al_down", dss.Al_down,
		"t_up", dss.T_up,
		"t_down", dss.T_down,
		"z_up", dss.Z_up,
		"z_down", dss.Z_down)
}
