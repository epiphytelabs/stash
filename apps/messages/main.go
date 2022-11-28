package main

import (
	"embed"
	_ "embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strings"

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

	c, err := stash.NewClient("api:4000")
	if err != nil {
		return err
	}

	g, err := graph.New(c)
	if err != nil {
		return err
	}

	v, err := vite()
	if err != nil {
		return err
	}

	s.Subrouter("/apps/messages", func(r *stdapi.Router) {
		r.Router.PathPrefix("/graph").Handler(g)
		r.Router.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			v.ServeHTTP(w, r)
		}))
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

	return viteProduction()
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

//go:embed dist/web/*
var dist embed.FS

func viteProduction() (http.Handler, error) {
	root, err := fs.Sub(dist, "dist/web")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return http.StripPrefix("/apps/messages", viteSinglePage(root)), nil
}

func viteSinglePage(dist fs.FS) http.Handler {
	hfs := http.FileServer(http.FS(dist))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := dist.Open(strings.TrimPrefix(r.URL.Path, "/")); errors.Is(err, fs.ErrNotExist) {
			r.URL.Path = "/"
		}
		hfs.ServeHTTP(w, r)
	})
}
