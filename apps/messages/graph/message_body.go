package graph

import "github.com/epiphytelabs/stash/apps/messages/pkg/message"

type MessageBody struct {
	body *message.Body
}

func (mb MessageBody) HTML() string {
	return mb.body.HTML
}

func (mb MessageBody) Text() string {
	return mb.body.Text
}
