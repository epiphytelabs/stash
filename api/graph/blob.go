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
	return transactionReturn(context.Background(), g, func(gg *Graph) (*Blob, error) {
		b, err := gg.store.BlobCreate(strings.NewReader(args.Data))
		if err != nil {
			return nil, err
		}

		if args.Labels != nil {
			for _, l := range *args.Labels {
				for _, v := range l.Values {
					if err := gg.store.LabelCreate(b.Hash, l.Key, v); err != nil {
						return nil, err
					}
				}
			}
		}

		return &Blob{*b, g}, nil
	})
}

type BlobRemovedArgs struct {
	Query string
}

func (g *Graph) BlobRemoved(ctx context.Context, args BlobAddedArgs) chan graphql.ID {
	ch := make(chan graphql.ID)

	sch := g.store.BlobRemoved(ctx, args.Query)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case id := <-sch:
				ch <- graphql.ID(id)
			}
		}
	}()

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

	lsh := map[string][]string{}

	for _, l := range ls {
		if _, ok := lsh[l.Key]; !ok {
			lsh[l.Key] = []string{}
		}

		lsh[l.Key] = append(lsh[l.Key], l.Value)
	}

	rls := []*Label{}

	for k, vs := range lsh {
		rls = append(rls, &Label{
			key:    k,
			values: vs,
		})
	}

	return rls, nil
}
