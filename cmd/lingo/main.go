package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/snivilised/li18ngo/internal/lo"
	"github.com/snivilised/li18ngo/tools/mio"
	"github.com/snivilised/li18ngo/tools/sorter"
)

func sort(transport mio.Transport, hashed bool) error {
	data, err := transport.Read()
	if err != nil {
		return errors.Wrapf(err, "error reading from transport: '%s'",
			transport.Address(),
		)
	}

	sorted, err := lo.TernaryE(hashed,
		func() ([]byte, error) {
			return sorter.Apply[sorter.HashedMessageEntry](data)
		},
		func() ([]byte, error) {
			return sorter.Apply[sorter.MessageEntry](data)
		},
	)

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
	path := flag.String("path", "",
		"Path to the input JSON file",
	)
	hashed := flag.Bool("hash", false,
		"switch to indicate if the translation file contains hashes",
	)

	flag.Parse()

	if *path == "" {
		fmt.Println("Usage: go run sort_messages.go -path <path-to-input-file> [-hash]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := sort(&mio.NativeReaderWriterFS{
		Path: *path,
	}, *hashed); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
