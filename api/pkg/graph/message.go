package graph

import (
	"github.com/epiphytelabs/stash/api/pkg/model"
)

type Message struct {
	model.Message
}

// func (g *Graph) Messages(ctx context.Context) ([]*Message, error) {
// 	id, err := g.user(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	search := store.Labels{
// 		"label[domain]": {"message"},
// 		"label[to]":     {id},
// 	}

// 	bs, err := g.store.BlobList(search)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ms := []*Message{}

// 	for _, b := range bs {
// 		r, err := g.store.BlobGet(b.Hash)
// 		if err != nil {
// 			return nil, err
// 		}

// 		m, err := message.New(r, b.Created)
// 		if err != nil {
// 			return nil, err
// 		}

// 		ms = append(ms, &Message{m})
// 	}

// 	return ms, nil
// }

// type MessageArgs struct {
// 	Hash graphql.ID
// }

// func (g *Graph) Message(args MessageArgs) (*Message, error) {
// 	b, err := g.store.BlobMetadata(string(args.Hash))
// 	if err != nil {
// 		return nil, err
// 	}

// 	r, err := g.store.BlobGet(string(args.Hash))
// 	if err != nil {
// 		return nil, err
// 	}

// 	m, err := message.New(r, b.Created)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Message{m}, nil
// }

func (m Message) Body() (*MessageBody, error) {
	b, err := m.Message.Body()
	if err != nil {
		return nil, err
	}

	return &MessageBody{b}, nil
}

func (m Message) From() (*MessageAddress, error) {
	a, err := m.Message.From()
	if err != nil {
		return nil, err
	}

	return &MessageAddress{*a}, nil
}

func (m Message) Received() DateTime {
	return DateTime{m.Message.Received}
}

func (m Message) Subject() (string, error) {
	return m.Message.Subject()
}
