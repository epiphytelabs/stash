package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	stash "github.com/epiphytelabs/stash/api/client"
	"github.com/epiphytelabs/stash/apps/messages/pkg/message"
)

const (
	filter = "email/v2"
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

	added := c.BlobAdded(context.Background(), unprocessedQuery())

	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()

	for {
		select {
		case b := <-added:
			if err := label(c, b); err != nil {
				log.Printf("error: %v\n", err)
			}
		case <-tick.C:
			bs, err := c.BlobList(unprocessedQuery())
			if err != nil {
				log.Printf("error: %v\n", err)
			}

			for _, b := range bs {
				if err := label(c, b); err != nil {
					log.Printf("error: %v\n", err)
				}
			}
		}
	}
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
		{Key: "filter", Values: []string{filter}},
		{Key: "thread", Values: []string{m.Thread()}},
	}

	if from := m.Header("From"); from != "" {
		if b.Labels.GetOne("from") != from {
			ls = append(ls, stash.Label{Key: "from", Values: []string{from}})
		}
	}

	if err := c.LabelAdd(b.Hash, ls); err != nil {
		return err
	}

	return nil
}

func unprocessedQuery() string {
	return fmt.Sprintf(`domain="message" filter!=%q`, filter)
}
