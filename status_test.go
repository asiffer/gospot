// spot_test.go

package gospot

import (
	"bytes"
	"encoding/json"
	"math"
	"os"
	"testing"
)

func TestInitStatus(t *testing.T) {
	title("Status")
}

func TestMarshallWithNaN(t *testing.T) {
	title("Marshalling with NaN")
	var out bytes.Buffer
	status := *NewSpotStatus()
	if js, err := json.Marshal(status); err != nil {
		t.Error(err)
	} else {
		json.Indent(&out, js, "", "    ")
		out.Write([]byte{'\n'})
		out.WriteTo(os.Stdout)
	}

	dstatus := DSpotStatus{
		Mean: 500.7,
		SpotStatus: SpotStatus{
			N:      50,
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
		},
	}
	if js, err := json.Marshal(dstatus); err != nil {
		t.Error(err)
	} else {
		json.Indent(&out, js, "", "    ")
		out.Write([]byte{'\n'})
		out.WriteTo(os.Stdout)
	}
}
