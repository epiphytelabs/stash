package store

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"sync"

	"github.com/ddollar/stdapi"
	"github.com/epiphytelabs/stash/api/pkg/root"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

type Store struct {
	db *pg.DB
	fs root.FS

	added   sync.Map
	removed sync.Map
}

func New(base string) (*Store, error) {
	if err := os.MkdirAll(base, 0755); err != nil {
		return nil, errors.WithStack(err)
	}

	opts, err := pg.ParseURL(os.Getenv("POSTGRES_URL"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	db := pg.Connect(opts)

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
		debug.PrintStack()
		return stdapi.Errorf(http.StatusBadRequest, "invalid hash")
	}

	return nil
}
