package api

import (
	"github.com/ddollar/stdapi"
	stash "github.com/epiphytelabs/stash/api/client"
	"github.com/pkg/errors"
)

func Run() error {
	s := stdapi.New("api", "api")

	c, err := stash.NewClient("api:4000")
	if err != nil {
		return err
	}

	g, err := NewGraph(c)
	if err != nil {
		return err
	}

	s.Router.Handle("/domains/mail/graph", g)

	if err := s.Listen("https", ":4000"); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
