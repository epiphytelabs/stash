package main

import (
	"fmt"
	"os"

	"github.com/epiphytelabs/stash/api/pkg/api"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	a, err := api.New("./db")
	if err != nil {
		return err
	}

	if err := a.Listen("https", ":4000"); err != nil {
		return errors.WithStack(err)
	}

	if err := a.Close(); err != nil {
		return err
	}

	return nil
}
