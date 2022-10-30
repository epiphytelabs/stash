package client

import (
	"io"

	"github.com/ddollar/stdsdk"
	"github.com/epiphytelabs/stash/api/pkg/store"
	"github.com/pkg/errors"
)

func (c *Client) BlobCreate(r io.Reader) (*store.Blob, error) {
	var b store.Blob

	opts := stdsdk.RequestOptions{Body: r}

	if err := c.Post("/blobs", opts, &b); err != nil {
		return nil, errors.WithStack(err)
	}

	return &b, nil
}

func (c *Client) BlobList() ([]store.Blob, error) {
	var bs []store.Blob

	if err := c.Get("/blobs", stdsdk.RequestOptions{}, &bs); err != nil {
		return nil, errors.WithStack(err)
	}

	return bs, nil
}
