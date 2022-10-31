package graph

import "github.com/graph-gophers/graphql-go"

type User struct {
	id string
}

func (g *Graph) User() (*User, error) {
	return nil, nil
	// return &User{id: "foo"}, nil
}

func (u User) ID() graphql.ID {
	return graphql.ID(u.id)
}
