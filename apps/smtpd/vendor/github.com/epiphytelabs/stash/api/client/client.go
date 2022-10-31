package client

import (
	"crypto/tls"
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

type Client struct {
	graphql graphql.Client
}

func NewClient(endpoint string) (*Client, error) {
	gc := graphql.NewClient(endpoint, &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})

	return &Client{graphql: gc}, nil
}

//go:generate go run github.com/Khan/genqlient genqlient.yml
