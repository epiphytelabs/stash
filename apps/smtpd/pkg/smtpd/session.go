package smtpd

import (
	"io"
	"log"

	"github.com/emersion/go-smtp"
	stash "github.com/epiphytelabs/stash/api/client"
)

type Session struct {
	stash *stash.Client
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

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	labels := stash.Labels{
		{Key: "domain", Values: []string{"message", "email"}},
		{Key: "from", Values: []string{s.from}},
		{Key: "to", Values: []string{s.to}},
	}

	// labels := map[string][]string{
	// 	"domain": {"message", "email"},
	// 	"from":   {s.from},
	// 	"to":     {s.to},
	// }

	b, err := s.stash.BlobCreate(string(data), labels)
	if err != nil {
		return err
	}

	log.Printf("ns=smtpd at=store from=%q to=%q hash=%q\n", s.from, s.to, b.Hash)

	// if err := s.stash.LabelCreate(b.Hash, labels); err != nil {
	// 	return err
	// }

	return nil
}

func (s *Session) Reset() {
	s.from = ""
	s.to = ""
}

func (s *Session) Logout() error {
	return nil
}
