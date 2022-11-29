package api

import (
	"context"
	"fmt"
	"time"

	"github.com/epiphytelabs/stash/api/pkg/coalesce"
	"github.com/graph-gophers/graphql-go"
)

type Thread struct {
	id string
	g  *Graph
}

func (g *Graph) Threads(ctx context.Context) ([]*Thread, error) {
	id, err := g.user(ctx)
	if err != nil {
		return nil, err
	}

	bs, err := g.stash.BlobList(fmt.Sprintf(`domain="mail" to=%q`, id))
	if err != nil {
		return nil, err
	}

	tsh := map[string]bool{}

	for _, b := range bs {
		tsh[coalesce.String(b.Labels.GetOne("thread"), b.Hash)] = true
	}

	ts := []*Thread{}

	for tid := range tsh {
		ts = append(ts, &Thread{tid, g})
	}

	return ts, nil
}

func (g *Graph) ThreadAdded(ctx context.Context) (chan *Thread, error) {
	id, err := g.user(ctx)
	if err != nil {
		return nil, err
	}

	cctx, cancel := context.WithCancel(ctx)

	bch := g.stash.BlobAdded(cctx, fmt.Sprintf(`domain="mail" to=%q`, id))

	ch := make(chan *Thread)

	go func() {
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				return
			case b := <-bch:
				fmt.Printf("b: %+v\n", b)
			}
		}
	}()

	return ch, nil
}

func (t Thread) ID() graphql.ID {
	return graphql.ID(t.id)
}

func (t Thread) Messages() ([]*Message, error) {
	bs, err := t.g.stash.BlobList(fmt.Sprintf(`domain="mail" thread=%q`, t.id))
	if err != nil {
		return nil, err
	}

	ms := []*Message{}

	for _, b := range bs {
		ms = append(ms, &Message{b, t.g})
	}

	return ms, nil
}

func (t Thread) Updated() (DateTime, error) {
	bs, err := t.g.stash.BlobList(fmt.Sprintf(`domain="mail" thread=%q`, t.id))
	if err != nil {
		return DateTime{}, err
	}

	var latest time.Time

	for _, b := range bs {
		if b.Created.After(latest) {
			latest = b.Created
		}
	}

	return DateTime{latest}, nil
}
