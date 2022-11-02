package graph

import (
	"net/mail"

	stash "github.com/epiphytelabs/stash/api/client"
	"github.com/epiphytelabs/stash/apps/messages/pkg/message"
	"github.com/pkg/errors"
)

type Message struct {
	blob stash.Blob
	g    *Graph
}

func (m Message) Body() (*MessageBody, error) {
	msg, err := m.msg()
	if err != nil {
		return nil, err
	}

	b, err := msg.Body()
	if err != nil {
		return nil, err
	}

	return &MessageBody{b}, nil
}

func (m Message) From() (*MessageAddress, error) {
	a, err := mail.ParseAddress(m.blob.Labels.GetOne("from"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &MessageAddress{*a}, nil
}

func (m Message) Received() DateTime {
	return DateTime{m.blob.Created}
}

func (m Message) Subject() (string, error) {
	msg, err := m.msg()
	if err != nil {
		return "", err
	}

	s, err := msg.Subject()
	if err != nil {
		return "", errors.WithStack(err)
	}

	return s, nil
}

func (m Message) msg() (*message.Message, error) {
	data, err := m.g.stash.BlobData(m.blob.Hash)
	if err != nil {
		return nil, err
	}

	msg, err := message.New(data)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return msg, nil
}
