package lab

import "fmt"

type InvalidByte struct {
	Offset int
	Value  byte
}

func (e *InvalidByte) Error() string {
	return fmt.Sprintf("byte 0x%02X at offset %d has its high bit set", e.Value, e.Offset)
}

type Receiver struct {
	Name    string   `json:"receiver"`
	Break   string   `json:"break"`
	Records []string `json:"records"`
	Tail    string   `json:"tail"`
}

type Report struct {
	Input     string     `json:"input"`
	Receivers []Receiver `json:"receivers"`
}

// Inspect checks the seven-bit byte rule, then applies the two example breaks
// from RFC 20. It retains any bytes after the final break as an incomplete tail.
func Inspect(input []byte) (Report, error) {
	for offset, value := range input {
		if value >= 0x80 {
			return Report{}, &InvalidByte{Offset: offset, Value: value}
		}
	}
	return Report{
		Input: string(input),
		Receivers: []Receiver{
			frame(input, "SRI", '.'),
			frame(input, "UCLA", '\r'),
		},
	}, nil
}

func frame(input []byte, name string, delimiter byte) Receiver {
	receiver := Receiver{Name: name, Break: string(delimiter), Records: []string{}}
	start := 0
	for offset, value := range input {
		if value == delimiter {
			receiver.Records = append(receiver.Records, string(input[start:offset]))
			start = offset + 1
		}
	}
	receiver.Tail = string(input[start:])
	return receiver
}
