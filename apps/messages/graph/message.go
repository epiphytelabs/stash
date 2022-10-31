package graph

import (
	"net/mail"
	"time"

	stash "github.com/epiphytelabs/stash/api/client"
	"github.com/epiphytelabs/stash/apps/messages/pkg/message"
	"github.com/pkg/errors"
)

type Message struct {
	labels   stash.Labels
	msg      *message.Message
	received time.Time
}

func (m Message) Body() (*MessageBody, error) {
	b, err := m.msg.Body()
	if err != nil {
		return nil, err
	}

	return &MessageBody{b}, nil
}

func (m Message) From() (*MessageAddress, error) {
	a, err := m.msg.From()
	if err != nil {
		return nil, err
	}

	if a != nil {
		return &MessageAddress{*a}, nil
	}

	from := m.labels.Get("from")

	if len(from) > 0 {
		la, err := mail.ParseAddress(from[0])
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return &MessageAddress{*la}, nil
	}

	return nil, errors.New("no from address")
}

func (m Message) Received() DateTime {
	return DateTime{m.received}
}

func (m Message) Subject() (string, error) {
	return m.msg.Subject()
}
