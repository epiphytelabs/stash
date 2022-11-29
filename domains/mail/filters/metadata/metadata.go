package metadata

import (
	"context"
	"fmt"
	"time"

	"github.com/ddollar/logger"
	stash "github.com/epiphytelabs/stash/api/client"
	"github.com/epiphytelabs/stash/api/pkg/coalesce"
	"github.com/epiphytelabs/stash/domains/mail/pkg/message"
)

const (
	filter = "mail/metadata/v3"
)

var (
	log = logger.New("ns=metadata")
)

func Run() error {
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
				log.Error(err) //nolint:errcheck
			}
		case <-tick.C:
			bs, err := c.BlobList(unprocessedQuery())
			if err != nil {
				log.Error(err) //nolint:errcheck
			}

			for _, b := range bs {
				if err := label(c, b); err != nil {
					log.Error(err) //nolint:errcheck
				}
			}
		}
	}
}

func label(c *stash.Client, b stash.Blob) error {
	log.At("label").Logf("hash=%q", b.Hash)

	data, err := c.BlobData(b.Hash)
	if err != nil {
		return err
	}

	m, err := message.New(data)
	if err != nil {
		return err
	}

	from := coalesce.String(b.Labels.GetOne("from"), b.Labels.GetOne("smtp/from"))
	to := coalesce.String(b.Labels.GetOne("to"), b.Labels.GetOne("smtp/to"))

	if fa, err := m.From(); fa != nil && err != nil {
		from = fa.String()
	}

	if ta, err := m.To(); ta != nil && err != nil {
		to = ta.Address
	}

	ls := stash.Labels{
		{Key: "filter", Values: []string{filter}},
		{Key: "from", Values: []string{from}},
		{Key: "thread", Values: []string{m.Thread()}},
		{Key: "to", Values: []string{to}},
	}

	if b.Labels.GetOne("smtp/to") == "" {
		ls = append(ls, stash.Label{Key: "smtp/to", Values: []string{b.Labels.GetOne("to")}})
	}

	if err := c.LabelAdd(b.Hash, ls); err != nil {
		return err
	}

	return nil
}

func unprocessedQuery() string {
	return fmt.Sprintf(`domain="mail" filter!=%q`, filter)
}
