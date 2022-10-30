package rest_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/ddollar/stdsdk"
	"github.com/epiphytelabs/stash/api/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenList(t *testing.T) {
	testServer(t, func(h *helpers) {
		opts := stdsdk.RequestOptions{
			Body: strings.NewReader("bar baz qux \n test this baz\n\n bar  foo\tqux\n\nbaz"),
		}

		var b store.Blob

		err := h.client.Post("/blobs", opts, &b)
		assert.NoError(t, err)
		require.NotNil(t, b)
		assert.Equal(t, "fb4182ff7a53d89b3953b2c155b8123fec6316f4e60b5f96f51187fead4451e2", b.Hash)

		var tokens map[string]int

		err = h.client.Get(fmt.Sprintf("/blobs/%s/tokens", b.Hash), stdsdk.RequestOptions{}, &tokens)
		require.NoError(t, err)
		assert.NoError(t, err)
		require.Len(t, tokens, 6)
		assert.Equal(t, 2, tokens["bar"])
		assert.Equal(t, 3, tokens["baz"])
		assert.Equal(t, 1, tokens["foo"])
		assert.Equal(t, 2, tokens["qux"])
		assert.Equal(t, 1, tokens["test"])
		assert.Equal(t, 1, tokens["this"])
	})
}

func TestTokenListNonexistant(t *testing.T) {
	testServer(t, func(h *helpers) {
		res, err := h.client.GetStream("/blobs/fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9/tokens", stdsdk.RequestOptions{})
		assert.EqualError(t, err, "hash not found: fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9")
		require.NotNil(t, res)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}
