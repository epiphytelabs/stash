package store

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/ddollar/stdapi"
	"github.com/epiphytelabs/stash/api/pkg/root"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type Store struct {
	db *sql.DB
	fs root.FS

	added   sync.Map
	removed sync.Map
}

func New(base string) (*Store, error) {
	if err := os.MkdirAll(base, 0755); err != nil {
		return nil, errors.WithStack(err)
	}

	db, err := initializeDatabase(base)
	if err != nil {
		return nil, err
	}

	s := &Store{
		db: db,
		fs: root.FS(base),
	}

	return s, nil
}

func (s *Store) Close() error {
	if err := s.db.Close(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Store) subscribeAdd(ctx context.Context, query string, ch chan Blob) {
	log.Printf("subscribing add: %v\n", query)
	s.added.Store(ch, query)
	<-ctx.Done()
	log.Printf("unsubscribing add: %v\n", query)
	s.added.Delete(ch)
}

func (s *Store) subscribeRemove(ctx context.Context, query string, ch chan string) {
	log.Printf("subscribing remove: %v\n", query)
	s.removed.Store(ch, query)
	<-ctx.Done()
	log.Printf("unsubscribing remove: %v\n", query)
	s.removed.Delete(ch)
}

func hashFile(hash string) string {
	return filepath.Join("hash", hash)
}

func hashValidate(hash string) error {
	if len(hash) != 64 {
		return stdapi.Errorf(http.StatusBadRequest, "invalid hash")
	}

	return nil
}

func initializeDatabase(base string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath.Join(base, "index.db"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := db.Ping(); err != nil {
		return nil, errors.WithStack(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS blobs (
			hash VARCHAR(64) PRIMARY KEY,
			size INTEGER NOT NULL,
			created TIMESTAMP NOT NULL DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW'))
		);

		CREATE TABLE IF NOT EXISTS labels (
			hash VARCHAR(64) NOT NULL,
			key VARCHAR NOT NULL,
			value VARCHAR NOT NULL
		);

		CREATE INDEX IF NOT EXISTS labels_hash ON labels (hash);
		CREATE INDEX IF NOT EXISTS labels_hash_key ON labels (hash, key);
		CREATE UNIQUE INDEX IF NOT EXISTS labels_hash_key_value ON labels (hash, key, value);
	`)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return db, nil
}
