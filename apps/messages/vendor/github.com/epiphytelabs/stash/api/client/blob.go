package client

import (
	"context"
	"io"
	"strings"
	"time"
)

type Blob struct {
	Created time.Time
	Hash    string
	Labels  Labels
}

func (c *Client) BlobData(id string) (io.Reader, error) {
	res, err := blobData(context.Background(), c.graphql, id)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(res.GetBlob().Data), nil
}

func (c *Client) BlobGet(id string) (*Blob, error) {
	res, err := blobGet(context.Background(), c.graphql, id)
	if err != nil {
		return nil, err
	}

	rb := res.GetBlob()

	b := Blob{
		Created: rb.Created,
		Hash:    rb.Hash,
		Labels:  Labels{},
	}

	for _, label := range rb.Labels {
		b.Labels = append(b.Labels, Label(label))
	}

	return &b, nil
}

func (c *Client) BlobList(labels []Label) ([]Blob, error) {
	ls := []LabelInput{}

	for _, label := range labels {
		ls = append(ls, LabelInput(label))
	}

	res, err := blobList(context.Background(), c.graphql, ls)
	if err != nil {
		return nil, err
	}

	bs := []Blob{}

	for _, rb := range res.GetBlobs() {
		b := Blob{
			Created: rb.Created,
			Hash:    rb.Hash,
		}

		for _, label := range rb.Labels {
			b.Labels = append(b.Labels, Label(label))
		}

		bs = append(bs, b)
	}

	return bs, nil
}
