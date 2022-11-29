package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/ddollar/migrate"
	"github.com/pkg/errors"
)

//go:embed migrations/*.sql
var migrations embed.FS

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	}
}

func run() error {
	if err := migrate.Run(os.Getenv("POSTGRES_URL"), migrations); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
