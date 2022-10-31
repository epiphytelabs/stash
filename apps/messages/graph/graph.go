package graph

import (
	"bufio"
	"context"
	_ "embed" // embed
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"

	stash "github.com/epiphytelabs/stash/api/client"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
	"github.com/pkg/errors"
)

type contextKey string

var contextClientCertificate = contextKey("client-certificate")

//go:embed schema.graphql
var schema string

type Graph struct {
	handler http.Handler
	stash   *stash.Client
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func New(c *stash.Client) (*Graph, error) {
	g := &Graph{
		stash: c,
	}

	schema, err := graphql.ParseSchema(schema, g, graphql.ErrorExtensioner(errorTracer))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	opts := graphqlws.WithContextGenerator(g)
	g.handler = graphqlws.NewHandlerFunc(schema, &relay.Handler{Schema: schema}, opts) // support http fallback

	return g, nil
}

func (g *Graph) BuildContext(ctx context.Context, r *http.Request) (context.Context, error) {
	info, err := url.QueryUnescape(r.Header.Get("X-Forwarded-Tls-Client-Cert-Info"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return context.WithValue(ctx, contextClientCertificate, info), nil
}

func (g *Graph) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := g.handler.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}

	c, rw, err := h.Hijack()
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return c, rw, nil
}

func (g *Graph) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Origin")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	switch r.Method {
	case "GET", "POST":
		g.handler.ServeHTTP(w, r)
	case "OPTIONS":
		fmt.Fprintf(w, "ok\n")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

var subjectParser = regexp.MustCompile(`^Subject="CN=(.*)"$`)

func (g *Graph) user(ctx context.Context) (string, error) {
	info, ok := ctx.Value(contextClientCertificate).(string)
	if !ok {
		return "", errors.Errorf("no client certificate")
	}

	m := subjectParser.FindStringSubmatch(info)

	if len(m) > 1 {
		return m[1], nil
	}

	return "", errors.Errorf("no user found in certificate")
}

func errorTracer(err error) map[string]interface{} {
	if os.Getenv("MODE") == "development" {
		if st, ok := err.(stackTracer); ok {
			return map[string]interface{}{
				"stacktrace": st.StackTrace(),
			}
		}
	}

	return nil
}
