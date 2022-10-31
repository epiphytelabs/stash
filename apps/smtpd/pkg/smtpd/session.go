package smtpd

import (
	"io"
	"log"

	"github.com/emersion/go-smtp"
	"github.com/epiphytelabs/stash/api/client"
)

type Session struct {
	stash *client.Client
	from  string
	to    string
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
	log.Printf("ns=smtpd at=data from=%q to=%q\n", s.from, s.to)

	b, err := s.stash.BlobCreate(r)
	if err != nil {
		return err
	}

	log.Printf("ns=smtpd at=store from=%q to=%q hash=%q\n", s.from, s.to, b.Hash)

	labels := map[string][]string{
		"domain": {"message", "email"},
		"from":   {s.from},
		"to":     {s.to},
	}

	if err := s.stash.LabelCreate(b.Hash, labels); err != nil {
		return err
	}

	return nil
}

func (s *Session) Reset() {
	s.from = ""
	s.to = ""
}

func (s *Session) Logout() error {
	return nil
}
