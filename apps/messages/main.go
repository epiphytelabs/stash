package main

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"

	"github.com/convox/stdapi"
	stash "github.com/epiphytelabs/stash/api/client"
	"github.com/epiphytelabs/stash/apps/messages/graph"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	s := stdapi.New("messages", "messages")

	c, err := stash.NewClient("https://api:4000/graph")
	if err != nil {
		return err
	}

	g, err := graph.New(c)
	if err != nil {
		return err
	}

	h, err := vite()
	if err != nil {
		return err
	}

	s.Subrouter("/apps/messages", func(r *stdapi.Router) {
		r.Router.PathPrefix("/graph").Handler(g)
		r.Router.PathPrefix("/").Handler(h)
	})

	if err := s.Listen("https", ":4000"); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func npmRebuild() error {
	cmd := exec.Command("npm", "rebuild")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func vite() (http.Handler, error) {
	if os.Getenv("MODE") == "development" {
		return viteDevelopment()
	}

	return nil, nil
}

func viteDevelopment() (http.Handler, error) {
	if err := npmRebuild(); err != nil {
		return nil, err
	}

	cmd := exec.Command("npx", "vite", "--host")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	go io.Copy(os.Stderr, stdout) //nolint:errcheck

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	go io.Copy(os.Stderr, stderr) //nolint:errcheck

	if err := cmd.Start(); err != nil {
		return nil, errors.WithStack(err)
	}

	u, err := url.Parse("http://localhost:3000/")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	h := httputil.NewSingleHostReverseProxy(u)

	return h, nil
}
