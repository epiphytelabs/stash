package rest

import (
	"net/http"

	"github.com/ddollar/stdapi"
	"github.com/epiphytelabs/stash/api/pkg/store"
	"github.com/pkg/errors"
)

func (r *REST) LabelCreate(ctx *stdapi.Context) error {
	hash := ctx.Var("hash")

	req := ctx.Request()

	if err := req.ParseForm(); err != nil {
		return errors.WithStack(err)
	}

	err := r.store.LabelCreate(hash, req.Form)
	switch err {
	case store.ErrHashNotFound:
		return stdapi.Errorf(http.StatusNotFound, "hash not found: %s", hash)
	case nil:
		return errors.WithStack(ctx.RenderOK())
	default:
		return err
	}
}

func (r *REST) LabelDelete(ctx *stdapi.Context) error {
	hash := ctx.Var("hash")

	err := r.store.LabelDelete(hash, ctx.Request().URL.Query())
	switch err {
	case store.ErrHashNotFound:
		return stdapi.Errorf(http.StatusNotFound, "hash not found: %s", hash)
	case nil:
		return errors.WithStack(ctx.RenderOK())
	default:
		return err
	}
}

func (r *REST) LabelGet(ctx *stdapi.Context) error {
	hash := ctx.Var("hash")
	key := ctx.Var("key")

	values, err := r.store.LabelGet(hash, key)
	switch err {
	case store.ErrHashNotFound:
		return stdapi.Errorf(http.StatusNotFound, "hash not found: %s", hash)
	case nil:
	default:
		return err
	}

	return errors.WithStack(ctx.RenderJSON(values))
}

func (r *REST) LabelList(ctx *stdapi.Context) error {
	hash := ctx.Var("hash")

	labels, err := r.store.LabelList(hash)
	switch err {
	case store.ErrHashNotFound:
		return stdapi.Errorf(http.StatusNotFound, "hash not found: %s", hash)
	case nil:
	default:
		return err
	}

	return errors.WithStack(ctx.RenderJSON(labels))
}
