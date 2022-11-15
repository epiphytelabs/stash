package graph

import (
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
	for _, l := range args.Labels {
		for _, v := range l.Values {
			if err := g.store.LabelCreate(string(args.Hash), l.Key, v); err != nil {
				return "", err
			}
		}
	}

	return args.Hash, nil
}
