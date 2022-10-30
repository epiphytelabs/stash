package graph

import (
	"github.com/epiphytelabs/stash/api/pkg/store"
	"github.com/graph-gophers/graphql-go"
)

type Blob struct {
	store.Blob
}

func (g *Graph) Blobs() ([]*Blob, error) {
	blobs, err := g.store.BlobList(nil)
	if err != nil {
		return nil, err
	}

	var bs []*Blob
	for _, b := range blobs {
		bs = append(bs, &Blob{b})
	}

	return bs, nil
}

func (b *Blob) Hash() graphql.ID {
	return graphql.ID(b.Blob.Hash)
}

func (b *Blob) Created() DateTime {
	return DateTime{b.Blob.Created}
}
