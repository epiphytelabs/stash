package smtpd

import (
	"time"

	"github.com/ddollar/logger"
	"github.com/emersion/go-smtp"
	stash "github.com/epiphytelabs/stash/api/client"
	"github.com/pkg/errors"
)

type Server struct {
	log   *logger.Logger
	stash *stash.Client
}

func NewServer(url string) (*Server, error) {
	c, err := stash.NewClient(url)
	if err != nil {
		return nil, err
	}

	return NewServerWithStash(c)
}

func NewServerWithStash(stash *stash.Client) (*Server, error) {
	s := &Server{
		log:   logger.New("ns=smtpd"),
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

	s.log.At("listen").Logf("addr=%q", addr)

	return errors.WithStack(ss.ListenAndServe())
}

func (s *Server) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return &Session{server: s}, nil
}

func (s *Server) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	return &Session{server: s}, nil
}
