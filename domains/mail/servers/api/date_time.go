package api

import (
	"time"

	"github.com/pkg/errors"
)

type DateTime struct {
	time.Time
}

func (DateTime) ImplementsGraphQLType(name string) bool {
	return name == "DateTime"
}

func (ts *DateTime) UnmarshalGraphQL(input interface{}) error {
	switch t := input.(type) {
	case string:
		tt, err := time.Parse("2006-01-02T15:04:05Z", t)
		if err != nil {
			return errors.WithStack(err)
		}

		ts.Time = tt
	default:
		return errors.Errorf("unknown DateTime unmarshal type: %T", t)
	}

	return nil
}
