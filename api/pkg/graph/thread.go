package graph

import (
	"context"

	"github.com/epiphytelabs/stash/api/pkg/model"
	"github.com/graph-gophers/graphql-go"
)

type Thread struct {
	model.Thread
}

func (g *Graph) Threads(ctx context.Context) ([]*Thread, error) {
	id, err := g.user(ctx)
	if err != nil {
		return nil, err
	}

	ts, err := g.model.ThreadList(id)
	if err != nil {
		return nil, err
	}

	gts := []*Thread{}

	for _, t := range ts {
		gts = append(gts, &Thread{t})
	}

	return gts, nil
}

func (t Thread) ID() graphql.ID {
	return graphql.ID(t.Thread.ID)
}

func (t Thread) Messages() []*Message {
	ms := []*Message{}

	for _, m := range t.Thread.Messages {
		ms = append(ms, &Message{m})
	}

	return ms
}

func (t Thread) Participants() ([]*MessageAddress, error) {
	ras := []*MessageAddress{}

	for _, a := range t.Thread.Participants {
		ras = append(ras, &MessageAddress{a})
	}

	return ras, nil
}

func (t Thread) Subject() string {
	return t.Thread.Subject
}

func (t Thread) Updated() DateTime {
	return DateTime{t.Thread.Updated}
}
