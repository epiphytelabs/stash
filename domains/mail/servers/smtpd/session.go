package smtpd

import (
	"bytes"
	"io"

	"github.com/emersion/go-smtp"
	stash "github.com/epiphytelabs/stash/api/client"
	"github.com/epiphytelabs/stash/domains/mail/pkg/message"
	"github.com/pkg/errors"
)

type Session struct {
	server *Server
	from   string
	to     string
}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	s.to = to
	return nil
}

func (s *Session) Data(r io.Reader) error {
	log := s.server.log.At("session").Append("from=%q to=%q", s.from, s.to)

	data, err := io.ReadAll(r)
	if err != nil {
		return errors.WithStack(err)
	}

	m, err := message.New(bytes.NewReader(data))
	if err != nil {
		return errors.WithStack(err)
	}

	from, err := m.From()
	if err != nil {
		return errors.WithStack(err)
	}
	if from == nil {
		return errors.New("missing from address")
	}

	labels := stash.Labels{
		{Key: "domain", Values: []string{"mail"}},
		{Key: "smtp/from", Values: []string{s.from}},
		{Key: "smtp/to", Values: []string{s.to}},
	}

	b, err := s.server.stash.BlobCreate(string(data), labels)
	if err != nil {
		return err
	}

	return log.Successf("hash=%q", b.Hash) //nolint:wrapcheck
}

func (s *Session) Reset() {
	s.from = ""
	s.to = ""
}

func (s *Session) Logout() error {
	return nil
}
