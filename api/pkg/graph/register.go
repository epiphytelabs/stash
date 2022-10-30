package graph

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"

	"github.com/epiphytelabs/stash/api/pkg/ca"
	"github.com/pkg/errors"
	"software.sslmate.com/src/go-pkcs12"
)

type RegisterArgs struct {
	ID       string
	Password string
}

func (g *Graph) Register(args RegisterArgs) (string, error) {
	c, err := ca.LoadFiles("/db/certs/public/ca.pem", "/db/certs/private/ca.pem")
	if err != nil {
		return "", err
	}

	cc, err := c.GenerateClient(args.ID)
	if err != nil {
		return "", err
	}

	cert, err := x509.ParseCertificate(cc.Certificate[0])
	if err != nil {
		return "", errors.WithStack(err)
	}

	data, err := pkcs12.Encode(rand.Reader, cc.PrivateKey, cert, nil, args.Password)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return base64.StdEncoding.EncodeToString(data), nil
}
