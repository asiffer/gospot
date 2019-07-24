// status.go

package gospot

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
)

// SpotStatus is the structure embedding the status of a Spot instance
type SpotStatus struct {
	// N is the number of normal observations (not the alarms)
	N int `json:"n"`
	// ExUp is the current number of up excesses
	ExUp int `json:"ex_up"`
	// ExDown is the current number of down excesses
	ExDown int `json:"ex_down"`
	// NtUp is the total number of up excesses
	NtUp int `json:"Nt_up"`
	// NtDown is the total number of down excesses
	NtDown int `json:"Nt_down"`
	// AlUp is the number of up alarms
	AlUp int `json:"al_up"`
	// AlDown is the number of down alarms
	AlDown int `json:"al_down"`
	// TUp is the transitional up threshold
	TUp float64 `json:"t_up"`
	// TDown is the transitional down threshold
	TDown float64 `json:"t_down"`
	// ZUp is the up alert thresholds
	ZUp float64 `json:"th_up"`
	// ZDown is the down alert thresholds
	ZDown float64 `json:"th_down"`
}

// NewSpotStatus creates a new empty spot status structure
func NewSpotStatus() *SpotStatus {
	return &SpotStatus{
		N:      0,
		ExUp:   0,
		ExDown: 0,
		NtUp:   0,
		NtDown: 0,
		AlUp:   0,
		AlDown: 0,
		TUp:    math.NaN(),
		TDown:  math.NaN(),
		ZUp:    math.NaN(),
		ZDown:  math.NaN(),
	}
}

// DSpotStatus is the structure embedding the status of a DSpot instance
type DSpotStatus struct {
	SpotStatus
	// Mean is the the value of the current local model
	Mean float64 `json:"drift"`
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

// PreMarshalWithNaN prepares to marshal a structure sending a map
// removing the NaN field/value
func PreMarshalWithNaN(x interface{}) map[string]interface{} {
	// pointer to struct - addressable
	v := reflect.ValueOf(x)
	t := v.Type()

	m := make(map[string]interface{})
	// var field string

	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i).Interface()

		field, exists := t.Field(i).Tag.Lookup("json")
		if !exists {
			field = t.Field(i).Name
		}

		switch value.(type) {
		case float64:
			if f, ok := value.(float64); !math.IsNaN(f) && ok {
				m[field] = value
			}
		default:
			m[field] = value
		}
	}
	return m
}

// MarshalJSON is the method required to implement the Marshaler interface.
// Marshaler is the interface implemented by types that can marshal
// themselves into valid JSON.
func (ss SpotStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(PreMarshalWithNaN(ss))
}

// MarshalJSON is the method required to implement the Marshaler interface.
// Marshaler is the interface implemented by types that can marshal
// themselves into valid JSON.
func (dss DSpotStatus) MarshalJSON() ([]byte, error) {
	m := PreMarshalWithNaN(dss.SpotStatus)
	m["drift"] = dss.Mean
	return json.Marshal(m)
}
