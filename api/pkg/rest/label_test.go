package rest_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ddollar/stdsdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLabelCreate(t *testing.T) {
	testServer(t, func(h *helpers) {
		b, err := h.store.BlobCreate(strings.NewReader("bar"))
		require.NoError(t, err)
		require.NotEmpty(t, b.Hash)

		opts := stdsdk.RequestOptions{
			Params: stdsdk.Params{
				"key1": []string{"val1a", "val1b"},
				"key2": "val2",
			},
		}

		err = h.client.Post(fmt.Sprintf("/blobs/%s/labels", b.Hash), opts, nil)
		require.NoError(t, err)

		labels, err := h.store.LabelList(b.Hash)
		assert.NoError(t, err)
		require.Len(t, labels, 2)
		require.Len(t, labels["key1"], 2)
		assert.Equal(t, "val1a", labels["key1"][0])
		assert.Equal(t, "val1b", labels["key1"][1])
		require.Len(t, labels["key2"], 1)
		assert.Equal(t, "val2", labels["key2"][0])
	})
}

func TestLabelCreateInvalidHash(t *testing.T) {
	testServer(t, func(h *helpers) {
		err := h.client.Post("/blobs/invalid/labels", stdsdk.RequestOptions{}, nil)
		assert.EqualError(t, err, "invalid hash")
	})
}

func TestLabelCreateInvalidLabel(t *testing.T) {
	testServer(t, func(h *helpers) {
		b, err := h.store.BlobCreate(strings.NewReader("bar"))
		require.NoError(t, err)
		require.NotEmpty(t, b.Hash)

		opts := stdsdk.RequestOptions{
			Params: stdsdk.Params{
				"": []string{"val1"},
			},
		}

		err = h.client.Post(fmt.Sprintf("/blobs/%s/labels", b.Hash), opts, nil)
		require.EqualError(t, err, "key required")
	})
}

func TestLabelCreateMissingBlob(t *testing.T) {
	testServer(t, func(h *helpers) {
		err := h.client.Post("/blobs/0000000011111111222222223333333344444444555555556666666677777777/labels", stdsdk.RequestOptions{}, nil)
		assert.EqualError(t, err, "hash not found: 0000000011111111222222223333333344444444555555556666666677777777")
	})
}

func TestLabelDelete(t *testing.T) {
	testServer(t, func(h *helpers) {
		b, err := h.store.BlobCreate(strings.NewReader("bar"))
		require.NoError(t, err)
		require.NotEmpty(t, b.Hash)

		require.NoError(t, h.store.LabelCreate(b.Hash, map[string][]string{
			"key1": {"val1a", "val1b"},
			"key2": {"val2"},
		}))

		opts := stdsdk.RequestOptions{
			Query: stdsdk.Query{
				"key1": []string{"val1a", "val1b"},
			},
		}

		err = h.client.Delete(fmt.Sprintf("/blobs/%s/labels", b.Hash), opts, nil)
		require.NoError(t, err)

		labels, err := h.store.LabelList(b.Hash)
		assert.NoError(t, err)
		require.Len(t, labels, 1)
		require.Len(t, labels["key2"], 1)
		assert.Equal(t, "val2", labels["key2"][0])
	})
}

func TestLabelDeleteInvalidLabel(t *testing.T) {
	testServer(t, func(h *helpers) {
		b, err := h.store.BlobCreate(strings.NewReader("bar"))
		require.NoError(t, err)
		require.NotEmpty(t, b.Hash)

		opts := stdsdk.RequestOptions{
			Query: stdsdk.Query{
				"": []string{"val1"},
			},
		}

		err = h.client.Delete(fmt.Sprintf("/blobs/%s/labels", b.Hash), opts, nil)
		require.EqualError(t, err, "key required")
	})
}

func TestLabelDeleteInvalidHash(t *testing.T) {
	testServer(t, func(h *helpers) {
		err := h.client.Delete("/blobs/invalid/labels", stdsdk.RequestOptions{}, nil)
		assert.EqualError(t, err, "invalid hash")
	})
}

func TestLabelDeleteMissingBlob(t *testing.T) {
	testServer(t, func(h *helpers) {
		err := h.client.Delete("/blobs/0000000011111111222222223333333344444444555555556666666677777777/labels", stdsdk.RequestOptions{}, nil)
		assert.EqualError(t, err, "hash not found: 0000000011111111222222223333333344444444555555556666666677777777")
	})
}

func TestLabelGet(t *testing.T) {
	testServer(t, func(h *helpers) {
		b, err := h.store.BlobCreate(strings.NewReader("bar"))
		require.NoError(t, err)
		require.NotEmpty(t, b.Hash)

		require.NoError(t, h.store.LabelCreate(b.Hash, map[string][]string{
			"key1": {"val1a", "val1b"},
		}))

		expected := []string{"val1a", "val1b"}

		var actual []string

		err = h.client.Get(fmt.Sprintf("/blobs/%s/labels/key1", b.Hash), stdsdk.RequestOptions{}, &actual)
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	})
}

func TestLabelGetInvalidLabel(t *testing.T) {
	testServer(t, func(h *helpers) {
		b, err := h.store.BlobCreate(strings.NewReader("bar"))
		require.NoError(t, err)
		require.NotEmpty(t, b.Hash)

		err = h.client.Get(fmt.Sprintf("/blobs/%s/labels/", b.Hash), stdsdk.RequestOptions{}, nil)
		require.EqualError(t, err, "key required")
	})
}

func TestLabelGetInvalidHash(t *testing.T) {
	testServer(t, func(h *helpers) {
		err := h.client.Get("/blobs/invalid/labels/key1", stdsdk.RequestOptions{}, nil)
		assert.EqualError(t, err, "invalid hash")
	})
}

func TestLabelGetMissingBlob(t *testing.T) {
	testServer(t, func(h *helpers) {
		err := h.client.Get("/blobs/0000000011111111222222223333333344444444555555556666666677777777/labels/key1", stdsdk.RequestOptions{}, nil)
		assert.EqualError(t, err, "hash not found: 0000000011111111222222223333333344444444555555556666666677777777")
	})
}

func TestLabelList(t *testing.T) {
	testServer(t, func(h *helpers) {
		b, err := h.store.BlobCreate(strings.NewReader("bar"))
		require.NoError(t, err)
		require.NotEmpty(t, b.Hash)

		require.NoError(t, h.store.LabelCreate(b.Hash, map[string][]string{
			"key1": {"val1a", "val1b"},
			"key2": {"val2"},
		}))

		opts := stdsdk.RequestOptions{
			Params: stdsdk.Params{
				"key1": []string{"val1a", "val1b"},
				"key2": "val2",
			},
		}

		var actual map[string][]string

		err = h.client.Get(fmt.Sprintf("/blobs/%s/labels", b.Hash), opts, &actual)
		require.NoError(t, err)

		assert.Len(t, actual, 2)
		assert.Equal(t, []string{"val1a", "val1b"}, actual["key1"])
		assert.Equal(t, []string{"val2"}, actual["key2"])
	})
}

func TestLabelListInvalidHash(t *testing.T) {
	testServer(t, func(h *helpers) {
		err := h.client.Get("/blobs/invalid/labels", stdsdk.RequestOptions{}, nil)
		assert.EqualError(t, err, "invalid hash")
	})
}

func TestLabelListMissingBlob(t *testing.T) {
	testServer(t, func(h *helpers) {
		err := h.client.Get("/blobs/0000000011111111222222223333333344444444555555556666666677777777/labels", stdsdk.RequestOptions{}, nil)
		assert.EqualError(t, err, "hash not found: 0000000011111111222222223333333344444444555555556666666677777777")
	})
}
