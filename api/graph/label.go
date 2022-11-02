package graph

import (
	"github.com/epiphytelabs/stash/api/internal/store"
	"github.com/graph-gophers/graphql-go"
)

type Label struct {
	key    string
	values []string
}

func (l *Label) Key() string {
	return l.key
}

func (l *Label) Values() []string {
	return l.values
}

type LabelAddArgs struct {
	Hash   graphql.ID
	Labels []struct {
		Key    string
		Values []string
	}
}

func (g *Graph) LabelAdd(args LabelAddArgs) (graphql.ID, error) {
	labels := store.Labels{}

	for _, l := range args.Labels {
		labels[l.Key] = append(labels[l.Key], l.Values...)
	}

	if err := g.store.LabelCreate(string(args.Hash), labels); err != nil {
		return "", err
	}

	return args.Hash, nil
}
