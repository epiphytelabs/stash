package client

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/pkg/errors"
)

type Blob struct {
	Created time.Time
	Hash    string
	Labels  Labels
}

type Subscription struct {
	Body any
}

func (c *Client) BlobAdded(ctx context.Context, query string) <-chan Blob {
	type subscription struct {
		BlobAdded Blob `graphql:"blobAdded(query: $query)"`
	}

	// var res sub

	vars := map[string]interface{}{
		"query": query,
	}

	ch := make(chan Blob)

	go subscribe(ctx, c.sub, vars, func(res subscription) {
		ch <- res.BlobAdded
	})

	return ch
}

func (c *Client) BlobCreate(data string, labels Labels) (*Blob, error) {
	var res struct {
		BlobCreate Blob `graphql:"blobCreate(data:$data, labels:$labels)"`
	}

	vars := map[string]interface{}{
		"data":   data,
		"labels": labels.Input(),
	}

	if err := c.graphql.Mutate(context.Background(), &res, vars); err != nil {
		return nil, errors.WithStack(err)
	}

	return &res.BlobCreate, nil
}

func (c *Client) BlobData(id string) (io.Reader, error) {
	var res struct {
		Blob struct{ Data string } `graphql:"blob(id: $id)"`
	}

	vars := map[string]interface{}{
		"id": graphql.ID(id),
	}

	if err := c.graphql.Query(context.Background(), &res, vars); err != nil {
		return nil, errors.WithStack(err)
	}

	return strings.NewReader(res.Blob.Data), nil
}

func (c *Client) BlobList(query string) ([]Blob, error) {
	var res struct {
		Blobs []Blob `graphql:"blobs(query: $query)"`
	}

	vars := map[string]interface{}{
		"query": query,
	}

	if err := c.graphql.Query(context.Background(), &res, vars); err != nil {
		return nil, errors.WithStack(err)
	}

	return res.Blobs, nil
}

func (c *Client) BlobRemoved(ctx context.Context, query string) <-chan string {
	type subscription struct {
		BlobRemoved graphql.ID `graphql:"blobRemoved(query: $query)"`
	}

	// var res sub

	vars := map[string]interface{}{
		"query": query,
	}

	ch := make(chan string)

	go subscribe(ctx, c.sub, vars, func(res subscription) {
		ch <- string(res.BlobRemoved)
	})

	return ch
}
