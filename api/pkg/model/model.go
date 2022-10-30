package model

import "github.com/epiphytelabs/stash/api/pkg/store"

type Model struct {
	store *store.Store
}

func New(s *store.Store) (*Model, error) {
	return &Model{s}, nil
}
