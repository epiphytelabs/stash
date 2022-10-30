package store

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	queryLabel = regexp.MustCompile(`label\[([^\]]+)\]`)
)

type Blob struct {
	Hash    string    `json:"hash"`
	Created time.Time `json:"created"`
}

func (s *Store) BlobCreate(r io.Reader) (*Blob, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	file := hashFile(hash)

	if err := s.BlobExists(hash); err == nil {
		fmt.Printf("err: %+v\n", err)
		return nil, errors.Errorf("hash exists: %s", hash)
	}

	if _, err := s.fs.Stat(file); !os.IsNotExist(err) {
		fmt.Printf("err: %+v\n", err)
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

	for token, count := range tokenize(string(data)) {
		if _, err := s.db.Exec("INSERT INTO tokens (hash, token, count) VALUES (?, ?, ?)", hash, token, count); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return s.BlobMetadata(hash)
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

func (s *Store) BlobGet(hash string) (io.ReadCloser, error) {
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

func (s *Store) BlobList(labels map[string][]string) ([]Blob, error) {
	query := "SELECT DISTINCT blobs.hash, blobs.created from blobs"

	ql, args, err := blobLabelQuery(labels)
	if err != nil {
		return nil, err
	}

	query += ql
	query += " ORDER BY blobs.created DESC"

	rows, err := s.db.Query(query, args...)
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

func (s *Store) BlobMetadata(hash string) (*Blob, error) {
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

func (s *Store) BlobNew(labels map[string][]string, since time.Time) ([]Blob, error) {
	query := "SELECT DISTINCT blobs.hash, blobs.created from blobs"

	q, args, err := blobLabelQuery(labels)
	if err != nil {
		return nil, err
	}

	query += q

	query += "WHERE created > ?"
	args = append(args, since)

	rows, err := s.db.Query(query, args...)
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

func blobLabelQuery(labels map[string][]string) (string, []any, error) {
	args := []any{}
	idx := map[string]int{}
	query := ""

	for key, values := range labels {
		if key == "search" {
			if len(values) == 1 {
				for _, token := range strings.Fields(values[0]) {
					query += fmt.Sprintf(" INNER JOIN tokens AS tk%d ON tk%d.hash = blobs.hash AND LOWER(tk%d.token) LIKE ?", idx["token"], idx["token"], idx["token"])
					args = append(args, fmt.Sprintf("%%%s%%", strings.ToLower(token)))
					idx["token"]++
				}
			}
		} else if m := queryLabel.FindStringSubmatch(key); len(m) > 1 {
			for _, value := range values {
				query += fmt.Sprintf(" INNER JOIN labels AS lbl%d ON lbl%d.hash = blobs.hash AND lbl%d.key = ? AND lbl%d.value = ?", idx["label"], idx["label"], idx["label"], idx["label"])
				args = append(args, m[1], value)
				idx["label"]++
			}
		} else {
			return "", nil, errors.Errorf("invalid query: %s", key)
		}
	}

	return query, args, nil
}
