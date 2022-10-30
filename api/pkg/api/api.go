package api

import (
	"fmt"
	"net/http"

	"github.com/ddollar/stdapi"
	"github.com/epiphytelabs/stash/api/pkg/graph"
	"github.com/epiphytelabs/stash/api/pkg/rest"
	"github.com/epiphytelabs/stash/api/pkg/store"
)

type API struct {
	*stdapi.Server
	graph *graph.Graph
	rest  *rest.REST
}

func New(base string) (*API, error) {
	s, err := store.New(base)
	if err != nil {
		return nil, err
	}

	return NewWithStore(s)
}

func NewWithStore(s *store.Store) (*API, error) {
	g, err := graph.New(s)
	if err != nil {
		return nil, err
	}

	r, err := rest.New(s)
	if err != nil {
		return nil, err
	}

	a := &API{
		Server: stdapi.New("stash", "stash"),
		graph:  g,
		rest:   r,
	}

	a.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	// pass, ok := os.LookupEnv("API_PASSWORD")
	// if !ok {
	// 	return nil, fmt.Errorf("API_PASSWORD not set")
	// }

	// auth := authenticate(pass)

	// a.Router.PathPrefix("/graphql").Handler(auth(g))
	// a.Router.PathPrefix("/").Handler(auth(r))

	a.Router.PathPrefix("/api").Handler(r)
	a.Router.PathPrefix("/graph").Handler(g)

	return a, nil
}

func (a *API) Close() error {
	if err := a.rest.Close(); err != nil {
		return err
	}

	return nil
}

// func authenticate(password string) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			if _, pass, _ := r.BasicAuth(); pass != password {
// 				w.WriteHeader(http.StatusUnauthorized)
// 				fmt.Fprintf(w, "unauthorized\n")
// 				return
// 			}

// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// func extractClientCertificate(next http.Handler) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		cert := r.Header.Get("X-Forwarded-Tls-Client-Cert")
// 		ctx := context.WithValue(r.Context(), "client-certificate", cert)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	}
// }
