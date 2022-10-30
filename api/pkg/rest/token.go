package rest

import (
	"net/http"

	"github.com/ddollar/stdapi"
	"github.com/epiphytelabs/stash/api/pkg/store"
	"github.com/pkg/errors"
)

func (r *REST) TokenList(ctx *stdapi.Context) error {
	hash := ctx.Var("hash")

	labels, err := r.store.TokenList(hash)
	switch err {
	case store.ErrHashNotFound:
		return stdapi.Errorf(http.StatusNotFound, "hash not found: %s", hash)
	case nil:
	default:
		return err
	}

	return errors.WithStack(ctx.RenderJSON(labels))
}
