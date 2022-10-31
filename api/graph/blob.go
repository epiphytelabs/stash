package graph

import (
	"io"

	"github.com/epiphytelabs/stash/api/pkg/store"
	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"
)

type Blob struct {
	blob store.Blob
	g    *Graph
}

type BlobArgs struct {
	ID graphql.ID
}

func (g *Graph) Blob(args BlobArgs) (*Blob, error) {
	b, err := g.store.BlobGet(string(args.ID))
	if err != nil {
		return nil, err
	}

	return &Blob{*b, g}, nil
}

type BlobsArgs struct {
	Labels *[]struct {
		Key    string
		Values []string
	}
}

func (g *Graph) Blobs(args BlobsArgs) ([]*Blob, error) {
	labels := store.Labels{}

	if args.Labels != nil {
		for _, l := range *args.Labels {
			labels[l.Key] = append(labels[l.Key], l.Values...)
		}
	}

	blobs, err := g.store.BlobList(labels)
	if err != nil {
		return nil, err
	}

	var bs []*Blob
	for _, b := range blobs {
		bs = append(bs, &Blob{b, g})
	}

	return bs, nil
}

func (b *Blob) Created() DateTime {
	return DateTime{b.blob.Created}
}

func (b *Blob) Data() (string, error) {
	r, err := b.g.store.BlobData(b.blob.Hash)
	if err != nil {
		return "", err
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(data), nil
}

func (b *Blob) Hash() graphql.ID {
	return graphql.ID(b.blob.Hash)
}

func (b *Blob) Labels() ([]*Label, error) {
	ls, err := b.g.store.LabelList(b.blob.Hash)
	if err != nil {
		return nil, err
	}

	rls := []*Label{}

	for k, vs := range ls {
		rls = append(rls, &Label{
			key:    k,
			values: vs,
		})
	}

	return rls, nil
}
