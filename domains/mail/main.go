package main

import (
	"fmt"
	"os"

	"github.com/epiphytelabs/stash/domains/mail/filters/metadata"
	"github.com/epiphytelabs/stash/domains/mail/servers/api"
	"github.com/epiphytelabs/stash/domains/mail/servers/smtpd"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	g := new(errgroup.Group)

	g.Go(api.Run)
	g.Go(metadata.Run)
	g.Go(smtpd.Run)

	if err := g.Wait(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
