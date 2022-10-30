package rest

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"

	"github.com/ddollar/stdapi"
	"github.com/epiphytelabs/stash/api/pkg/ca"
	"github.com/pkg/errors"
	"software.sslmate.com/src/go-pkcs12"
)

func (r *REST) UserCreate(ctx *stdapi.Context) error {
	id := ctx.Form("id")
	password := ctx.Form("password")

	c, err := ca.LoadFiles("/db/certs/public/ca.pem", "/db/certs/private/ca.pem")
	if err != nil {
		return err
	}

	cc, err := c.GenerateClient(id)
	if err != nil {
		return err
	}

	cert, err := x509.ParseCertificate(cc.Certificate[0])
	if err != nil {
		return errors.WithStack(err)
	}

	data, err := pkcs12.Encode(rand.Reader, cc.PrivateKey, cert, nil, password)
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(ctx.RenderText(base64.StdEncoding.EncodeToString(data)))
}
