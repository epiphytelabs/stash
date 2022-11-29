package api

import (
	"time"

	"github.com/pkg/errors"
)

type Date struct {
	time.Time
}

func (Date) ImplementsGraphQLType(name string) bool {
	return name == "Date"
}

func (ts *Date) UnmarshalGraphQL(input interface{}) error {
	switch t := input.(type) {
	case string:
		tt, err := time.Parse("2006-01-02", t)
		if err != nil {
			return errors.WithStack(err)
		}

		ts.Time = tt
	default:
		return errors.Errorf("unknown Date unmarshal type: %T", t)
	}

	return nil
}
