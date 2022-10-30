package graph

import (
	"github.com/epiphytelabs/stash/api/pkg/model"
)

type Message struct {
	model.Message
}

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
