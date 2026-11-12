package lab

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestAllByteValues(t *testing.T) {
	for value := 0; value <= 255; value++ {
		_, err := Inspect([]byte{'A', byte(value)})
		if value < 128 {
			if err != nil {
				t.Errorf("valid ASCII byte 0x%02X: %v", value, err)
			}
			continue
		}
		var invalid *InvalidByte
		if !errors.As(err, &invalid) || invalid.Offset != 1 || invalid.Value != byte(value) {
			t.Errorf("byte 0x%02X: expected exact offset and value, got %v", value, err)
		}
	}
}

func TestReceiverExamples(t *testing.T) {
	got, err := Inspect([]byte("PAY 12.50\rACK.\r"))
	if err != nil {
		t.Fatal(err)
	}
	want := []Receiver{
		{Name: "SRI", Break: ".", Records: []string{"PAY 12", "50\rACK"}, Tail: "\r"},
		{Name: "UCLA", Break: "\r", Records: []string{"PAY 12.50", "ACK."}},
	}
	if !reflect.DeepEqual(got.Receivers, want) {
		t.Fatalf("receivers = %+v, want %+v", got.Receivers, want)
	}
}

func TestFramingPreservesBytes(t *testing.T) {
	for _, input := range []string{"", "..", "\r\r", "unterminated", "\x00.\r\x7f\n"} {
		report, err := Inspect([]byte(input))
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range report.Receivers {
			var rebuilt strings.Builder
			for _, record := range r.Records {
				rebuilt.WriteString(record)
				rebuilt.WriteString(r.Break)
			}
			rebuilt.WriteString(r.Tail)
			if rebuilt.String() != input {
				t.Errorf("%s reconstructed %q, want %q", r.Name, rebuilt.String(), input)
			}
		}
	}
	report, err := Inspect([]byte(".."))
	if err != nil || !reflect.DeepEqual(report.Receivers[0].Records, []string{"", ""}) {
		t.Fatalf("two breaks must produce two empty records: %+v, %v", report, err)
	}
}
