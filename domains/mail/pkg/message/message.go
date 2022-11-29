package message

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"

	"github.com/epiphytelabs/stash/api/pkg/clean"
	"github.com/epiphytelabs/stash/api/pkg/coalesce"
	"github.com/pkg/errors"
)

type Message struct {
	msg  *mail.Message
	hash string
}

type Messages []Message

type Body struct {
	HTML string
	Text string
}

func New(r io.Reader) (*Message, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	m := &Message{
		msg:  msg,
		hash: fmt.Sprintf("%x", sha256.Sum256(data)),
	}

	return m, nil
}

func (m *Message) Body() (*Body, error) {
	mth := m.msg.Header.Get("Content-Type")
	if mth == "" {
		data, err := io.ReadAll(m.msg.Body)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return &Body{Text: string(data)}, nil
	}

	mt, params, err := mime.ParseMediaType(mth)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	switch mt {
	case "text/html":
		data, err := io.ReadAll(quotedprintable.NewReader(m.msg.Body))
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return &Body{HTML: clean.Mail(string(data))}, nil
	case "text/plain":
		data, err := io.ReadAll(m.msg.Body)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		switch m.msg.Header.Get("Content-Transfer-Encoding") {
		case "base64":
			dec, err := base64.StdEncoding.DecodeString(string(data))
			if err != nil {
				return nil, errors.WithStack(err)
			}

			return &Body{Text: string(dec)}, nil
		default:
			return nil, errors.Errorf("unknown content transfer encoding: %s", m.msg.Header.Get("Content-Transfer-Encoding"))
		}
	default:
		if !strings.HasPrefix(mt, "multipart/") {
			return nil, errors.Errorf("not multipart: %s", mt)
		}
	}

	mr := multipart.NewReader(m.msg.Body, params["boundary"])
	mb := &Body{}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.WithStack(err)
		}

		data, err := io.ReadAll(part)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		pmt, _, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
		if err != nil {
			return nil, errors.WithStack(err)
		}

		switch pmt {
		case "text/html":
			mb.HTML += clean.Mail(string(data))
		case "text/plain":
			mb.Text += string(data)
		}
	}

	return mb, nil
}

func (m *Message) From() (*mail.Address, error) {
	if from := m.msg.Header.Get("From"); from != "" {
		a, err := mail.ParseAddress(from)
		if err != nil {
			return nil, errors.Wrapf(err, from)
		}

		return a, nil
	}

	return nil, nil
}

func (m *Message) Header(name string) string {
	return m.msg.Header.Get(name)
}

func (m *Message) Subject() (string, error) {
	s, err := new(mime.WordDecoder).DecodeHeader(m.msg.Header.Get("Subject"))
	if err != nil {
		return "", errors.WithStack(err)
	}

	return s, nil
}

func (m *Message) Thread() string {
	return coalesce.String(m.msg.Header.Get("In-Reply-To"), m.msg.Header.Get("Message-ID"), m.hash)
}

func (m *Message) To() (*mail.Address, error) {
	if to := m.msg.Header.Get("To"); to != "" {
		a, err := mail.ParseAddress(to)
		if err != nil {
			return nil, errors.Wrapf(err, to)
		}

		return a, nil
	}

	return nil, nil
}
