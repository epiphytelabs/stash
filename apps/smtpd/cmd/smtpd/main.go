package main

import (
	"fmt"
	"os"

	"github.com/epiphytelabs/stash/apps/smtpd/pkg/smtpd"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	s, err := smtpd.New(os.Getenv("STASH") + "/api")
	if err != nil {
		return err
	}

	return s.Listen(":2525")
}
