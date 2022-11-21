package store

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

func Transaction(ctx context.Context, s *Store, fn func(*Store) error) error {
	if _, ok := s.db.(*pg.Tx); ok {
		return fn(s)
	}

	return errors.WithStack(s.db.RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		return fn(&Store{db: tx, fs: s.fs})
	}))
}

func TransactionReturn[R any](ctx context.Context, s *Store, fn func(*Store) (R, error)) (R, error) {
	var ret R

	err := Transaction(ctx, s, func(s *Store) error {
		res, err := fn(s)
		if err != nil {
			return err
		}

		ret = res

		return nil
	})

	return ret, err
}
