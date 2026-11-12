package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"example.com/rfc-go-labs/lab"
)

func TestSampleOutputs(t *testing.T) {
	input, err := os.ReadFile("../../examples/interchange.dat")
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := run(nil, bytes.NewReader(input), &output); err != nil {
		t.Fatal(err)
	}
	var got lab.Report
	if err := json.Unmarshal(output.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Input != "PAY 12.50\rACK.\r" || len(got.Receivers) != 2 || got.Receivers[0].Tail != "\r" {
		t.Fatalf("unexpected sample report: %+v", got)
	}
	output.Reset()
	if err := run([]string{"-html"}, strings.NewReader("<script>\r"), &output); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output.String(), "<script>") || !strings.Contains(output.String(), "&lt;script&gt;") {
		t.Fatal("expected escaped input in HTML")
	}
}

func TestRejectedInputProducesNoReport(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"A.café", "byte 0xC3 at offset 5"},
		{"A.\x80", "byte 0x80 at offset 2"},
		{strings.Repeat("a", (64<<10)+1), "input exceeds 64 KiB"},
	}
	for _, tc := range cases {
		var output bytes.Buffer
		err := run(nil, strings.NewReader(tc.input), &output)
		if err == nil || !strings.Contains(err.Error(), tc.want) || output.Len() != 0 {
			t.Errorf("want %q and no report, got %v and %d output bytes", tc.want, err, output.Len())
		}
	}
}
