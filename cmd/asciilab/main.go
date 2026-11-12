package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/rfc-go-labs/lab"
)

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout); err != nil && !errors.Is(err, flag.ErrHelp) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, input io.Reader, output io.Writer) error {
	flags := flag.NewFlagSet("asciilab", flag.ContinueOnError)
	html := flags.Bool("html", false, "render a static HTML report")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("usage: asciilab [-html] < input-file")
	}
	const maxInput = 64 << 10
	data, err := io.ReadAll(io.LimitReader(input, maxInput+1))
	if err != nil {
		return err
	}
	if len(data) > maxInput {
		return fmt.Errorf("input exceeds 64 KiB")
	}
	report, err := lab.Inspect(data)
	if err != nil {
		return err
	}
	if *html {
		return page.Execute(output, report)
	}
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}
