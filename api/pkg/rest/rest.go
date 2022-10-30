package rest

import (
	"github.com/ddollar/stdapi"
	"github.com/epiphytelabs/stash/api/pkg/store"
)

type REST struct {
	*stdapi.Server
	store *store.Store
}

func New(s *store.Store) (*REST, error) {
	r := &REST{
		Server: stdapi.New("stash", "stash"),
		store:  s,
	}

	r.Subrouter("/api", func(sub *stdapi.Router) {
		r.Routes(sub)
	})

	return r, nil
}

func (r *REST) Close() error {
	if err := r.store.Close(); err != nil {
		return err
	}

	return nil
}
