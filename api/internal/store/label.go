package store

import (
	"net/http"

	"github.com/ddollar/stdapi"
	"github.com/pkg/errors"
)

type Labels = map[string][]string

func (s *Store) LabelCreate(hash string, labels Labels) error {
	if err := hashValidate(hash); err != nil {
		return err
	}

	if err := s.BlobExists(hash); err != nil {
		return err
	}

	for key, values := range labels {
		if err := labelValidate(key); err != nil {
			return err
		}

		for _, value := range values {
			if _, err := s.db.Exec("INSERT INTO labels (hash, key, value) VALUES (?, ?, ?)", hash, key, value); err != nil {
				return errors.WithStack(err)
			}
		}
	}

	return nil
}

func (s *Store) LabelDelete(hash string, labels Labels) error {
	if err := hashValidate(hash); err != nil {
		return err
	}

	if err := s.BlobExists(hash); err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return errors.WithStack(err)
	}
	defer tx.Rollback() //nolint:errcheck

	for key, values := range labels {
		if err := labelValidate(key); err != nil {
			return err
		}

		for _, value := range values {
			if _, err := s.db.Exec("DELETE FROM labels WHERE hash = ? AND key = ? AND value = ?", hash, key, value); err != nil {
				return errors.WithStack(err)
			}
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

	if err := s.BlobExists(hash); err != nil {
		return nil, err
	}

	rows, err := s.db.Query("SELECT value FROM labels WHERE hash = ? AND key = ? ORDER BY value", hash, key)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	values := []string{}

	for rows.Next() {
		var value string

		if err := rows.Scan(&value); err != nil {
			return nil, errors.WithStack(err)
		}

		values = append(values, value)
	}

	return values, nil
}

func (s *Store) LabelList(hash string) (Labels, error) {
	if err := hashValidate(hash); err != nil {
		return nil, err
	}

	if err := s.BlobExists(hash); err != nil {
		return nil, err
	}

	rows, err := s.db.Query("SELECT key, value FROM labels WHERE hash = ?", hash)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	labels := Labels{}

	for rows.Next() {
		var key, value string

		if err := rows.Scan(&key, &value); err != nil {
			return nil, errors.WithStack(err)
		}

		labels[key] = append(labels[key], value)
	}

	return labels, nil
}

func labelValidate(key string) error {
	if len(key) == 0 {
		return stdapi.Errorf(http.StatusBadRequest, "key required")
	}

	return nil
}
