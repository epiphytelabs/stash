package rest_test

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ddollar/logger"
	"github.com/ddollar/stdsdk"
	"github.com/epiphytelabs/stash/api/pkg/rest"
	"github.com/epiphytelabs/stash/api/pkg/store"
	"github.com/stretchr/testify/require"
)

type helpers struct {
	client *stdsdk.Client
	rest   *rest.REST
	store  *store.Store
}

func testServer(t *testing.T, fn func(*helpers)) {
	tmp, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	s, err := store.New(tmp)
	require.NoError(t, err)

	r, err := rest.New(s)
	require.NoError(t, err)

	r.Logger = logger.Discard
	r.Server.Recover = func(err error) {
		require.NoError(t, err, "httptest server panic")
	}

	ht := httptest.NewServer(r)
	defer ht.Close()

	c, err := stdsdk.New(ht.URL + "/api")
	require.NoError(t, err)

	h := &helpers{
		client: c,
		rest:   r,
		store:  s,
	}

	fn(h)

	err = os.RemoveAll(tmp)
	require.NoError(t, err)
}
