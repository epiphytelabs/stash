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
	Created time.Time `json:"created"`
}

func (s *Store) BlobAdded(ctx context.Context, query string) chan Blob {
	ch := make(chan Blob)
	go s.subscribeAdd(ctx, query, ch)
	return ch
}

func (s *Store) BlobCreate(r io.Reader) (*Blob, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	file := hashFile(hash)

	if err := s.BlobExists(hash); err == nil {
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
}

func (s *Store) BlobDelete(hash string) error {
	if err := hashValidate(hash); err != nil {
		return err
	}

	if err := s.BlobExists(hash); err != nil {
		return err
	}

	if _, err := s.db.Exec("DELETE FROM blobs WHERE hash = ?", hash); err != nil {
		return errors.WithStack(err)
	}

	file := hashFile(hash)

	exists, err := s.fs.Exists(file)
	if err != nil {
		return err
	} else if !exists {
		return ErrHashNotFound
	}

	if err := s.fs.Remove(file); err != nil {
		return err
	}

	return nil
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

func (s *Store) BlobExists(hash string) error {
	if err := hashValidate(hash); err != nil {
		return err
	}

	rows, err := s.db.Query("SELECT hash FROM blobs WHERE hash = ?", hash)
	if err != nil {
		return errors.WithStack(err)
	}
	defer rows.Close()

	if !rows.Next() {
		return ErrHashNotFound
	}

	return nil
}

func (s *Store) BlobGet(hash string) (*Blob, error) {
	rows, err := s.db.Query("SELECT created FROM blobs WHERE hash = ?", hash)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, ErrHashNotFound
	}

	blob := Blob{Hash: hash}

	if err := rows.Scan(&blob.Created); err != nil {
		return nil, errors.WithStack(err)
	}

	return &blob, nil
}

func (s *Store) BlobList(query string) ([]Blob, error) {
	q := "SELECT DISTINCT blobs.hash, blobs.created from blobs WHERE "

	qf, args, err := blobQueryFragment(query)
	if err != nil {
		return nil, err
	}

	q += qf
	q += " ORDER BY blobs.created DESC"

	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	blobs := []Blob{}

	for rows.Next() {
		var blob Blob
		if err := rows.Scan(&blob.Hash, &blob.Created); err != nil {
			return nil, errors.WithStack(err)
		}

		blobs = append(blobs, blob)
	}

	return blobs, nil
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

	rows, err := s.db.Query(q, args...)
	if err != nil {
		return false, errors.WithStack(err)
	}
	defer rows.Close()

	if !rows.Next() {
		return false, errors.New("no rows in result set")
	}

	var count int

	if err := rows.Scan(&count); err != nil {
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
