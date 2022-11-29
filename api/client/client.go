package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/pkg/errors"
)

type Client struct {
	graphql *graphql.Client
	sub     *graphql.SubscriptionClient
}

func NewClient(host string) (*Client, error) {
	hc := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 10 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       300 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	wsopts := graphql.WebsocketOptions{
		HTTPClient: hc,
	}

	c := &Client{
		graphql: graphql.NewClient(fmt.Sprintf("https://%s/graph", host), hc),
		sub:     graphql.NewSubscriptionClient(fmt.Sprintf("wss://%s/graph", host)).WithWebSocketOptions(wsopts),
		// sub:     graphql.NewSubscriptionClient(fmt.Sprintf("wss://%s/graph", host)).WithWebSocketOptions(wsopts).WithRetryTimeout(1 * time.Second).WithLog(log.Println),
	}

	c.sub.OnDisconnected(func() {
		// c.sub.Reset()
	})
	// c.sub.OnDisconnected(func() {
	// 	fmt.Println("reconnecting")
	// 	c.sub.Reset()
	// })

	c.sub.OnError(func(sc *graphql.SubscriptionClient, err error) error {
		if err != nil {
			sc.Close()
		}
		return sc.Reset() //nolint:wrapcheck
	})

	return c, nil
}

func subscribe[T any](ctx context.Context, c *graphql.SubscriptionClient, vars map[string]any, fn func(T)) {
	var query T

	id, err := c.Subscribe(&query, vars, func(data []byte, err error) error {
		if err := json.Unmarshal(data, &query); err != nil {
			return errors.WithStack(err)
		}

		fn(query)

		return nil
	})
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	go c.Run() //nolint:errcheck

	<-ctx.Done()

	if err := c.Unsubscribe(id); err != nil {
		log.Printf("error: %v\n", err)
	}
}
