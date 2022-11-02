package main

import (
	"context"
	"fmt"
	"log"
	"os"

	stash "github.com/epiphytelabs/stash/api/client"
	"github.com/epiphytelabs/stash/apps/messages/pkg/message"
)

const (
	filter = "email/v0"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	c, err := stash.NewClient("api:4000")
	if err != nil {
		return err
	}

	bs, err := c.BlobList(unprocessedQuery())
	if err != nil {
		return err
	}

	for _, b := range bs {
		if err := label(c, b); err != nil {
			log.Printf("error: %v\n", err)
		}
	}

	ch := c.BlobAdded(context.Background(), unprocessedQuery())

	for b := range ch {
		if err := label(c, b); err != nil {
			log.Printf("error: %v\n", err)
		}
	}

	return nil
}

func label(c *stash.Client, b stash.Blob) error {
	log.Printf("label hash=%s\n", b.Hash)

	data, err := c.BlobData(b.Hash)
	if err != nil {
		return err
	}

	m, err := message.New(data)
	if err != nil {
		return err
	}

	ls := stash.Labels{
		{Key: "filter", Values: []string{"email/v0"}},
		{Key: "thread", Values: []string{m.Thread()}},
	}

	if err := c.LabelAdd(b.Hash, ls); err != nil {
		return err
	}

	return nil
}

func unprocessedQuery() string {
	return fmt.Sprintf(`domain="message" filter!=%q`, filter)
}
