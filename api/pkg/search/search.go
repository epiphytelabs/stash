package search

import (
	"github.com/alecthomas/participle/v2"
	"github.com/pkg/errors"
)

type Query struct {
	Fields []*Field `parser:"@@*"`
}

type Field struct {
	Key   string `parser:"@Ident"`
	Op    string `parser:"@('=' | '!' '=')"`
	Value string `parser:"@String"`
}

var parser = participle.MustBuild[Query](participle.Unquote("String"))

func Parse(query string) (*Query, error) {
	q, err := parser.ParseString("", query)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return q, nil
}
