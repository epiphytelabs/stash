package client

import (
	"fmt"

	"github.com/ddollar/stdsdk"
	"github.com/epiphytelabs/stash/api/pkg/store"
	"github.com/pkg/errors"
)

func (c *Client) LabelCreate(hash string, labels store.Labels) error {
	params := stdsdk.Params{}

	for k := range labels {
		params[k] = labels[k]
	}

	opts := stdsdk.RequestOptions{Params: params}

	if err := c.Post(fmt.Sprintf("/blobs/%s/labels", hash), opts, nil); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
