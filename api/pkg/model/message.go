package model

import (
	"net/mail"
	"time"

	"github.com/epiphytelabs/stash/api/pkg/message"
	"github.com/epiphytelabs/stash/api/pkg/store"
	"github.com/pkg/errors"
)

type Message struct {
	*message.Message
	Received time.Time
	blob     store.Blob
	labels   store.Labels
}

func (m *Model) MessageList(to string) ([]Message, error) {
	search := store.Labels{
		"label[domain]": {"message"},
		"label[to]":     {to},
	}

	bs, err := m.store.BlobList(search)
	if err != nil {
		return nil, err
	}

	ms := []Message{}

	for _, b := range bs {
		r, err := m.store.BlobGet(b.Hash)
		if err != nil {
			return nil, err
		}

		mm, err := message.New(r)
		if err != nil {
			return nil, err
		}

		ls, err := m.store.LabelList(b.Hash)
		if err != nil {
			return nil, err
		}

		ms = append(ms, Message{
			Message:  mm,
			Received: b.Created,
			blob:     b,
			labels:   ls,
		})
	}

	return ms, nil
}

func (m *Message) From() (*mail.Address, error) {
	f, err := m.Message.From()
	if err != nil {
		return nil, err
	}
	if f != nil {
		return f, nil
	}

	if from := m.labels["from"]; len(from) == 1 {
		a, err := mail.ParseAddress(from[0])
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return a, nil
	}

	return nil, errors.New("no from address")
}
