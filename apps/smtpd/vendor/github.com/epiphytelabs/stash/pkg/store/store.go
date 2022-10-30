package store

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ddollar/stdapi"
	"github.com/epiphytelabs/stash/pkg/root"
	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
	fs root.FS
}

func New(base string) (*Store, error) {
	if err := os.MkdirAll(base, 0755); err != nil {
		return nil, err
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
		return err
	}

	return nil
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
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
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
		CREATE INDEX IF NOT EXISTS labels_hash_key_value ON labels (hash, key, value);

		CREATE TABLE IF NOT EXISTS tokens (
			hash VARCHAR(64) NOT NULL,
			token VARCHAR NOT NULL,
			count INTEGER NOT NULL
		);

		CREATE INDEX IF NOT EXISTS tokens_hash ON tokens (hash);
		CREATE INDEX IF NOT EXISTS tokens_hash_token ON tokens (hash, token);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
