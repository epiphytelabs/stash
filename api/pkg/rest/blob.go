package rest

import (
	"io"
	"net/http"

	"github.com/ddollar/stdapi"
	"github.com/epiphytelabs/stash/api/pkg/store"
	"github.com/pkg/errors"
)

func (r *REST) BlobCreate(ctx *stdapi.Context) error {
	hash, err := r.store.BlobCreate(ctx)
	if err != nil {
		return err
	}

	return errors.WithStack(ctx.RenderJSON(hash))
}

func (r *REST) BlobDelete(ctx *stdapi.Context) error {
	hash := ctx.Var("hash")

	err := r.store.BlobDelete(hash)
	switch err {
	case store.ErrHashNotFound:
		return stdapi.Errorf(http.StatusNotFound, "hash not found: %s", hash)
	case nil:
		return errors.WithStack(ctx.RenderOK())
	default:
		return err
	}
}

func (r *REST) BlobExists(ctx *stdapi.Context) error {
	hash := ctx.Var("hash")

	err := r.store.BlobExists(hash)
	switch err {
	case store.ErrHashNotFound:
		return stdapi.Errorf(http.StatusNotFound, "hash not found: %s", hash)
	case nil:
		return errors.WithStack(ctx.RenderOK())
	default:
		return err
	}
}

func (r *REST) BlobGet(ctx *stdapi.Context) error {
	hash := ctx.Var("hash")

	rd, err := r.store.BlobGet(hash)
	switch err {
	case store.ErrHashNotFound:
		return stdapi.Errorf(http.StatusNotFound, "hash not found: %s", hash)
	case nil:
	default:
		return err
	}
	defer rd.Close()

	if _, err := io.Copy(ctx, rd); err != nil {
		return errors.Errorf("read failed")
	}

	return nil
}

func (r *REST) BlobList(ctx *stdapi.Context) error {
	bs, err := r.store.BlobList(ctx.Request().URL.Query())
	if err != nil {
		return err
	}

	return errors.WithStack(ctx.RenderJSON(bs))
}
