package client

import (
	"context"
	"sort"

	"github.com/hasura/go-graphql-client"
	"github.com/pkg/errors"
)

type Label struct {
	Key    string
	Values []string
}

type Labels []Label

type LabelInput struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

func (c *Client) LabelAdd(hash string, labels Labels) error {
	var res struct {
		LabelAdd graphql.ID `graphql:"labelAdd(hash:$hash, labels:$labels)"`
	}

	vars := map[string]interface{}{
		"hash":   graphql.ID(hash),
		"labels": labels.Input(),
	}

	if err := c.graphql.Mutate(context.Background(), &res, vars); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (ls Labels) Get(key string) []string {
	var values []string

	for _, l := range ls {
		if l.Key == key {
			values = append(values, l.Values...)
		}
	}

	return values
}

func (ls Labels) GetOne(key string) string {
	vs := ls.Get(key)
	sort.Strings(vs)

	if len(vs) > 0 {
		return vs[0]
	}

	return ""
}

func (ls Labels) Input() []LabelInput {
	var lis []LabelInput

	for _, l := range ls {
		lis = append(lis, LabelInput(l))
	}

	return lis
}
