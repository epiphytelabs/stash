//go:build tools

package main

import (
	_ "github.com/cespare/reflex"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/goware/modvendor"
)
