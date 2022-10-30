package graph

import (
	"net/mail"

	"github.com/epiphytelabs/stash/api/pkg/coalesce"
)

type MessageAddress struct {
	address mail.Address
}

func (ma MessageAddress) Address() string {
	return ma.address.Address
}

func (ma MessageAddress) Display() string {
	return coalesce.String(ma.address.Name, ma.address.Address)
}

func (ma MessageAddress) Name() string {
	return ma.address.Name
}
