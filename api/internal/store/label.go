package store

import (
	"log"
	"net/http"

	"github.com/ddollar/stdapi"
	"github.com/pkg/errors"
)

type Label struct {
	Hash  string `pg:",pk"`
	Key   string `pg:",pk"`
	Value string
}

type Labels []Label

func (s *Store) LabelCreate(hash, key, value string) error {
	unmatched := map[chan Blob]bool{}

	// TODO use a custom map so we can handle errors better
	s.added.Range(func(ch, query any) bool {
		m, err := s.blobMatch(hash, query.(string))
		if err != nil {
			log.Printf("error: %v\n", err)
			return true
		}

		if !m {
			unmatched[ch.(chan Blob)] = true
		}

		return true
	})

	if err := hashValidate(hash); err != nil {
		return err
	}

	exists, err := s.BlobExists(hash)
	if err != nil {
		return err
	}
	if !exists {
		return ErrHashNotFound
	}

	l := Label{
		Hash:  hash,
		Key:   key,
		Value: value,
	}

	if _, err := s.db.Model(&l).Insert(); err != nil {
		return errors.WithStack(err)
	}

	b, err := s.BlobGet(hash)
	if err != nil {
		return err
	}

	s.added.Range(func(ch, query any) bool {
		m, err := s.blobMatch(hash, query.(string))
		if err != nil {
			log.Printf("error: %v\n", err)
			return true
		}

		if unmatched[ch.(chan Blob)] && m {
			ch.(chan Blob) <- *b
		}
		return true
	})

	return nil
}

func (s *Store) LabelDelete(hash string, labels Labels) error {
	if err := hashValidate(hash); err != nil {
		return err
	}

	exists, err := s.BlobExists(hash)
	if err != nil {
		return err
	}
	if !exists {
		return ErrHashNotFound
	}

	tx, err := s.db.Begin()
	if err != nil {
		return errors.WithStack(err)
	}
	defer tx.Rollback() //nolint:errcheck

	for _, l := range labels {
		if err := labelValidate(l.Key); err != nil {
			return err
		}

		if _, err := s.db.Model(&l).WherePK().Delete(); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Store) LabelGet(hash, key string) ([]string, error) {
	if err := hashValidate(hash); err != nil {
		return nil, err
	}

	if err := labelValidate(key); err != nil {
		return nil, err
	}

	exists, err := s.BlobExists(hash)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrHashNotFound
	}

	values := []string{}

	if _, err := s.db.Query(&values, "SELECT value FROM labels WHERE hash = ? AND key = ? ORDER BY value", hash, key); err != nil {
		return nil, errors.WithStack(err)
	}

	return values, nil
}

func (s *Store) LabelList(hash string) (Labels, error) {
	if err := hashValidate(hash); err != nil {
		return nil, err
	}

	exists, err := s.BlobExists(hash)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrHashNotFound
	}

	var ls Labels

	if _, err := s.db.Query(&ls, "SELECT * FROM labels WHERE hash = ?", hash); err != nil {
		return nil, errors.WithStack(err)
	}

	return ls, nil
}

func labelValidate(key string) error {
	if len(key) == 0 {
		return stdapi.Errorf(http.StatusBadRequest, "key required")
	}

	return nil
}
