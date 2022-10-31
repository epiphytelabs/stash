package graph

import (
	"context"
	"sort"
	"time"

	stash "github.com/epiphytelabs/stash/api/client"
	"github.com/epiphytelabs/stash/apps/messages/pkg/message"
	"github.com/graph-gophers/graphql-go"
)

type Thread struct {
	id       string
	messages []*Message
	updated  time.Time
}

func (g *Graph) Threads(ctx context.Context) ([]*Thread, error) {
	id, err := g.user(ctx)
	if err != nil {
		return nil, err
	}

	labels := []stash.Label{
		{Key: "domain", Values: []string{"message"}},
		{Key: "to", Values: []string{id}},
	}

	bs, err := g.stash.BlobList(labels)
	if err != nil {
		return nil, err
	}

	bsh := map[string]stash.Blob{}
	msh := map[string]*message.Message{}
	tsh := map[string][]string{}

	for _, b := range bs {
		br, err := g.stash.BlobData(b.Hash)
		if err != nil {
			return nil, err
		}

		m, err := message.New(br)
		if err != nil {
			return nil, err
		}

		tid := m.Thread()

		if _, ok := tsh[tid]; !ok {
			tsh[tid] = []string{}
		}

		bsh[b.Hash] = b
		msh[b.Hash] = m
		tsh[tid] = append(tsh[tid], b.Hash)
	}

	ts := []*Thread{}

	for tid := range tsh {
		t := &Thread{
			id:       tid,
			messages: []*Message{},
		}

		for _, mid := range tsh[tid] {
			t.messages = append(t.messages, &Message{
				labels:   bsh[mid].Labels,
				msg:      msh[mid],
				received: bsh[mid].Created,
			})

			if bsh[mid].Created.After(t.updated) {
				t.updated = bsh[mid].Created
			}
		}

		ts = append(ts, t)
	}

	sort.Slice(ts, func(i, j int) bool {
		return ts[i].updated.After(ts[j].updated)
	})

	return ts, nil
}

func (t Thread) ID() graphql.ID {
	return graphql.ID(t.id)
}

func (t Thread) Messages() []*Message {
	return t.messages
}

func (t Thread) Updated() DateTime {
	return DateTime{t.updated}
}
