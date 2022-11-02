package graph

import (
	"context"
	"io"
	"strings"

	"github.com/epiphytelabs/stash/api/internal/store"
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
	Query string
}

func (g *Graph) Blobs(args BlobsArgs) ([]*Blob, error) {
	blobs, err := g.store.BlobList(args.Query)
	if err != nil {
		return nil, err
	}

	var bs []*Blob
	for _, b := range blobs {
		bs = append(bs, &Blob{b, g})
	}

	return bs, nil
}

type BlobAddedArgs struct {
	Query string
}

func (g *Graph) BlobAdded(ctx context.Context, args BlobAddedArgs) chan *Blob {
	ch := make(chan *Blob)

	sch := g.store.BlobAdded(ctx, args.Query)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case b := <-sch:
				ch <- &Blob{b, g}
			}
		}
	}()

	return ch
}

type BlobCreateArgs struct {
	Data   string
	Labels *[]struct {
		Key    string
		Values []string
	}
}

func (g *Graph) BlobCreate(args BlobCreateArgs) (*Blob, error) {
	b, err := g.store.BlobCreate(strings.NewReader(args.Data))
	if err != nil {
		return nil, err
	}

	labels := store.Labels{}

	if args.Labels != nil {
		for _, l := range *args.Labels {
			labels[l.Key] = append(labels[l.Key], l.Values...)
		}
	}

	if err := g.store.LabelCreate(b.Hash, labels); err != nil {
		return nil, err
	}

	return &Blob{*b, g}, nil
}

type BlobRemovedArgs struct {
	Labels *[]struct {
		Key    string
		Values []string
	}
}

func (g *Graph) BlobRemoved(args BlobAddedArgs) chan graphql.ID {
	ch := make(chan graphql.ID)

	return ch
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
