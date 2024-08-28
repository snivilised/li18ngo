package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/snivilised/li18ngo/tools/mio"
	"github.com/snivilised/li18ngo/tools/sorter"
)

func process(transport mio.Transport) error {
	data, err := transport.Read()
	if err != nil {
		return errors.Wrapf(err, "error reading from transport: '%s'",
			transport.Address(),
		)
	}

	sorted, err := sorter.Sort(data)
	if err != nil {
		return fmt.Errorf("error sorting: %v", err)
	}

	if err := transport.Write(sorted); err != nil {
		return errors.Wrapf(err, "error writing to transport: '%v'",
			transport.Address(),
		)
	}

	return nil
}

func main() {
	path := flag.String("path", "", "Path to the input JSON file")
	flag.Parse()

	if *path == "" {
		fmt.Println("Usage: go run sort_messages.go -path <path-to-input-file>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := process(&mio.NativeReaderWriterFS{
		Path: *path,
	}); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
