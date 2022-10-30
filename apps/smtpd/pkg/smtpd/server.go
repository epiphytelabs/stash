package smtpd

import (
	"log"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/epiphytelabs/stash/api/client"
	"github.com/pkg/errors"
)

type Server struct {
	stash *client.Client
}

func New(url string) (*Server, error) {
	c, err := client.New(url)
	if err != nil {
		return nil, err
	}

	return NewWithStash(c)
}

func NewWithStash(stash *client.Client) (*Server, error) {
	s := &Server{
		stash: stash,
	}

	return s, nil
}

func (s *Server) Listen(addr string) error {
	ss := smtp.NewServer(s)

	ss.Addr = addr
	ss.ReadTimeout = 10 * time.Second
	ss.WriteTimeout = 10 * time.Second
	ss.MaxMessageBytes = 1024 * 1024
	ss.MaxRecipients = 50
	ss.AllowInsecureAuth = true

	log.Printf("ns=stash.smtpd at=listen addr=%q\n", addr)

	return errors.WithStack(ss.ListenAndServe())
}

func (s *Server) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return &Session{stash: s.stash}, nil
}

func (s *Server) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	return &Session{stash: s.stash}, nil
}
