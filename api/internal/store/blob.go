package store

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/epiphytelabs/stash/api/pkg/search"
	"github.com/pkg/errors"
)

type Blob struct {
	Hash    string    `json:"hash"`
	Size    int       `json:"size"`
	Created time.Time `json:"created"`
}

func (s *Store) BlobAdded(ctx context.Context, query string) chan Blob {
	ch := make(chan Blob)
	go s.subscribeAdd(ctx, query, ch)
	return ch
}

func (s *Store) BlobCreate(r io.Reader) (*Blob, error) {
	return TransactionReturn(context.Background(), s, func(s *Store) (*Blob, error) {
		data, err := io.ReadAll(r)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		hash := fmt.Sprintf("%x", sha256.Sum256(data))
		file := hashFile(hash)

		exists, err := s.BlobExists(hash)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.Errorf("hash exists: %s", hash)
		}

		if _, err := s.fs.Stat(file); !os.IsNotExist(err) {
			return nil, errors.Errorf("hash exists: %s", hash)
		}

		if _, err := s.db.Exec("INSERT INTO blobs (hash, size) VALUES (?, ?)", hash, len(data)); err != nil {
			return nil, errors.WithStack(err)
		}

		f, err := s.fs.Create(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		if _, err := f.Write(data); err != nil {
			return nil, errors.WithStack(err)
		}

		return s.BlobGet(hash)
	})
}

func (s *Store) BlobDelete(hash string) error {
	return Transaction(context.Background(), s, func(s *Store) error {
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

		if _, err := s.db.Exec("DELETE FROM blobs WHERE hash = ?", hash); err != nil {
			return errors.WithStack(err)
		}

		file := hashFile(hash)

		exists, err = s.fs.Exists(file)
		if err != nil {
			return err
		} else if !exists {
			return ErrHashNotFound
		}

		if err := s.fs.Remove(file); err != nil {
			return err
		}

		return nil
	})
}

func (s *Store) BlobData(hash string) (io.ReadCloser, error) {
	if err := hashValidate(hash); err != nil {
		return nil, err
	}

	file := hashFile(hash)

	f, err := s.fs.Open(file)
	if os.IsNotExist(err) {
		return nil, ErrHashNotFound
	} else if err != nil {
		return nil, err
	}

	return f, nil
}

func (s *Store) BlobExists(hash string) (bool, error) {
	if err := hashValidate(hash); err != nil {
		return false, err
	}

	var b Blob

	if _, err := s.db.Query(&b, "SELECT * FROM blobs WHERE hash = ?", hash); err != nil {
		return false, errors.WithStack(err)
	}

	return b.Hash != "", nil
}

func (s *Store) BlobGet(hash string) (*Blob, error) {
	var b Blob

	if _, err := s.db.Query(&b, "SELECT * FROM blobs WHERE hash = ?", hash); err != nil {
		return nil, errors.WithStack(err)
	}

	return &b, nil
}

func (s *Store) BlobList(query string) ([]Blob, error) {
	q := "SELECT DISTINCT blobs.hash, blobs.created from blobs WHERE "

	qf, args, err := blobQueryFragment(query)
	if err != nil {
		return nil, err
	}

	q += qf
	q += " ORDER BY blobs.created DESC"

	bs := []Blob{}

	if _, err := s.db.Query(&bs, q, args...); err != nil {
		return nil, errors.WithStack(err)
	}

	return bs, nil
}

func (s *Store) BlobRemoved(ctx context.Context, query string) chan string {
	ch := make(chan string)
	go s.subscribeRemove(ctx, query, ch)
	return ch
}

func (s *Store) blobMatch(hash string, query string) (bool, error) {
	qf, args, err := blobQueryFragment(query)
	if err != nil {
		return false, err
	}

	q := "SELECT COUNT(*) AS count FROM blobs WHERE " + qf

	q += " AND blobs.hash = ?"
	args = append(args, hash)

	var count int

	if _, err := s.db.Query(&count, q, args...); err != nil {
		return false, errors.WithStack(err)
	}

	return count > 0, nil
}

func blobQueryFragment(query string) (string, []any, error) {
	q, err := search.Parse(query)
	if err != nil {
		return "", nil, err
	}

	fragment := "1=1"
	args := []any{}

	for _, f := range q.Fields {
		sub := "SELECT hash FROM labels WHERE key = ? AND value = ?"

		switch f.Op {
		case "=":
			fragment += fmt.Sprintf(" AND hash IN (%s)", sub)
			args = append(args, f.Key, f.Value)
		case "!=":
			fragment += fmt.Sprintf(" AND hash NOT IN (%s)", sub)
			args = append(args, f.Key, f.Value)
		default:
			return "", nil, errors.Errorf("unknown op: %s", f.Op)
		}
	}

	return fragment, args, nil
}
