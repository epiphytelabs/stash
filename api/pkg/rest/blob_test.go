package rest_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ddollar/stdsdk"
	"github.com/epiphytelabs/stash/api/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlobCreate(t *testing.T) {
	testServer(t, func(h *helpers) {
		opts := stdsdk.RequestOptions{
			Body: strings.NewReader("bar"),
		}

		var b store.Blob

		err := h.client.Post("/blobs", opts, &b)
		assert.NoError(t, err)
		require.NotNil(t, b)
		assert.Equal(t, "fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9", b.Hash)

		r, err := h.store.BlobGet(b.Hash)
		assert.NoError(t, err)
		require.NotNil(t, r)
		defer r.Close()
		data, err := io.ReadAll(r)
		assert.NoError(t, err)
		assert.Equal(t, "bar", string(data))
	})
}

func TestBlobCreateDuplicate(t *testing.T) {
	testServer(t, func(h *helpers) {
		err := h.client.Post("/blobs", stdsdk.RequestOptions{Body: strings.NewReader("bar")}, nil)
		require.NoError(t, err)

		err2 := h.client.Post("/blobs", stdsdk.RequestOptions{Body: strings.NewReader("bar")}, nil)
		require.EqualError(t, err2, "hash exists: fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9")
	})
}

func TestBlobCreateDeleteRecreate(t *testing.T) {
	testServer(t, func(h *helpers) {
		var b store.Blob

		err := h.client.Post("/blobs", stdsdk.RequestOptions{Body: strings.NewReader("bar")}, &b)
		assert.NoError(t, err)
		require.NotNil(t, b)
		assert.Equal(t, "fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9", b.Hash)

		r, err := h.store.BlobGet(b.Hash)
		assert.NoError(t, err)
		require.NotNil(t, r)
		defer r.Close()
		data, err := io.ReadAll(r)
		assert.NoError(t, err)
		assert.Equal(t, "bar", string(data))

		err = h.client.Delete(fmt.Sprintf("/blobs/%s", b.Hash), stdsdk.RequestOptions{}, nil)
		require.NoError(t, err)

		err = h.client.Post("/blobs", stdsdk.RequestOptions{Body: strings.NewReader("bar")}, &b)
		assert.NoError(t, err)
		require.NotNil(t, b)
		assert.Equal(t, "fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9", b.Hash)

		r, err = h.store.BlobGet(b.Hash)
		assert.NoError(t, err)
		require.NotNil(t, r)
		defer r.Close()
		data, err = io.ReadAll(r)
		assert.NoError(t, err)
		assert.Equal(t, "bar", string(data))
	})
}

func TestBlobDelete(t *testing.T) {
	testServer(t, func(h *helpers) {
		b, err := h.store.BlobCreate(strings.NewReader("bar"))
		assert.NoError(t, err)
		require.NotEmpty(t, b.Hash)

		err = h.client.Delete(fmt.Sprintf("/blobs/%s", b.Hash), stdsdk.RequestOptions{}, nil)
		assert.Nil(t, err)

		err = h.store.BlobExists(b.Hash)
		assert.EqualError(t, err, "hash not found")
	})
}

func TestBlobDeleteInvalidHash(t *testing.T) {
	testServer(t, func(h *helpers) {
		err := h.client.Delete("/blobs/invalid", stdsdk.RequestOptions{}, nil)
		assert.EqualError(t, err, "invalid hash")
	})
}

func TestBlobDeleteNonExistant(t *testing.T) {
	testServer(t, func(h *helpers) {
		err := h.client.Delete("/blobs/fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9", stdsdk.RequestOptions{}, nil)
		assert.EqualError(t, err, "hash not found: fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9")
	})
}

func TestBlobExists(t *testing.T) {
	testServer(t, func(h *helpers) {
		b, err := h.store.BlobCreate(strings.NewReader("bar"))
		assert.NoError(t, err)
		require.NotNil(t, b)
		require.NotEmpty(t, b.Hash)

		var exists bool
		err = h.client.Head(fmt.Sprintf("/blobs/%s", b.Hash), stdsdk.RequestOptions{}, &exists)
		require.NoError(t, err)
		assert.True(t, exists)
	})
}

func TestBlobExistsInvalidHash(t *testing.T) {
	testServer(t, func(h *helpers) {
		err := h.client.Head("/blobs/invalid", stdsdk.RequestOptions{}, nil)
		assert.EqualError(t, err, "response status 400")
	})
}

func TestBlobExistsNonexistant(t *testing.T) {
	testServer(t, func(h *helpers) {
		req, err := h.client.Request("GET", "/blobs/fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9", stdsdk.RequestOptions{})
		require.NoError(t, err)

		res, err := h.client.HandleRequest(req)
		require.EqualError(t, err, "hash not found: fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9")
		require.NotNil(t, res)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)

		var exists bool
		err = h.client.Head("/blobs/fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9", stdsdk.RequestOptions{}, &exists)
		require.EqualError(t, err, "response status 404")
		assert.False(t, exists)
	})
}

func TestBlobGet(t *testing.T) {
	testServer(t, func(h *helpers) {
		b, err := h.store.BlobCreate(strings.NewReader("bar"))
		assert.NoError(t, err)
		require.NotNil(t, b)
		assert.NotEmpty(t, b.Hash)

		res, err := h.client.GetStream(fmt.Sprintf("/blobs/%s", b.Hash), stdsdk.RequestOptions{})
		require.NoError(t, err)
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, "bar", string(data))
	})
}

func TestBlobGetInvalidHash(t *testing.T) {
	testServer(t, func(h *helpers) {
		err := h.client.Get("/blobs/invalid", stdsdk.RequestOptions{}, nil)
		assert.EqualError(t, err, "invalid hash")
	})
}

func TestBlobGetNonexistant(t *testing.T) {
	testServer(t, func(h *helpers) {
		res, err := h.client.GetStream("/blobs/fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9", stdsdk.RequestOptions{})
		assert.EqualError(t, err, "hash not found: fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9")
		require.NotNil(t, res)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func TestBlobListLabels(t *testing.T) {
	testServer(t, func(h *helpers) {
		b1, err := h.store.BlobCreate(strings.NewReader("foo bar baz"))
		require.NoError(t, err)
		require.NotNil(t, b1)
		require.NoError(t, h.store.LabelCreate(b1.Hash, map[string][]string{"key": {"val1", "val2", "val3"}}))

		b2, err := h.store.BlobCreate(strings.NewReader("bar baz qux"))
		require.NoError(t, err)
		require.NotNil(t, b2)
		require.NoError(t, h.store.LabelCreate(b2.Hash, map[string][]string{"key": {"val2", "val3", "val4"}}))

		b3, err := h.store.BlobCreate(strings.NewReader("baz qux quux"))
		require.NoError(t, err)
		require.NotNil(t, b3)
		require.NoError(t, h.store.LabelCreate(b3.Hash, map[string][]string{"key": {"val3", "val4", "val5"}}))

		var actual []store.Blob

		expected1 := []store.Blob{*b2, *b1}
		opts1 := stdsdk.RequestOptions{Query: stdsdk.Query{"label[key]": "val2"}}
		err = h.client.Get("/blobs", opts1, &actual)
		assert.NoError(t, err)
		assert.Equal(t, expected1, actual)

		expected2 := []store.Blob{*b3, *b2}
		opts2 := stdsdk.RequestOptions{Query: stdsdk.Query{"label[key]": []string{"val3", "val4"}}}
		err = h.client.Get("/blobs", opts2, &actual)
		assert.NoError(t, err)
		assert.Equal(t, expected2, actual)

		expected3 := []store.Blob{}
		opts3 := stdsdk.RequestOptions{Query: stdsdk.Query{"label[other]": "foo"}}
		err = h.client.Get("/blobs", opts3, &actual)
		assert.NoError(t, err)
		assert.Equal(t, expected3, actual)

		expected4 := []store.Blob{*b2, *b1}
		opts4 := stdsdk.RequestOptions{Query: stdsdk.Query{"search": "bar"}}
		err = h.client.Get("/blobs", opts4, &actual)
		assert.NoError(t, err)
		assert.Equal(t, expected4, actual)

		expected5 := []store.Blob{*b2}
		opts5 := stdsdk.RequestOptions{Query: stdsdk.Query{"label[key]": "val4", "search": "bar"}}
		err = h.client.Get("/blobs", opts5, &actual)
		assert.NoError(t, err)
		assert.Equal(t, expected5, actual)

		opts6 := stdsdk.RequestOptions{Query: stdsdk.Query{"unknown": "foo"}}
		err = h.client.Get("/blobs", opts6, &actual)
		assert.EqualError(t, err, "invalid query: unknown")
	})
}
