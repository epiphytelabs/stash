package client

import (
	"github.com/ddollar/stdsdk"
	"github.com/pkg/errors"
)

type Client struct {
	*stdsdk.Client
}

func New(url string) (*Client, error) {
	c, err := stdsdk.New(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Client{c}, nil
}
