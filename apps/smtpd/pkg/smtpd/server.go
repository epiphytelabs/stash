package smtpd

import (
	"log"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/epiphytelabs/stash/pkg/store"
	"github.com/pkg/errors"
)

type Server struct {
	store *store.Store
}

func New(base string) (*Server, error) {
	s, err := store.New(base)
	if err != nil {
		return nil, err
	}

	return NewWithStore(s)
}

func NewWithStore(store *store.Store) (*Server, error) {
	s := &Server{
		store: store,
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
	return &Session{store: s.store}, nil
}

func (s *Server) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	return &Session{store: s.store}, nil
}
