package stdapi

import (
	"fmt"
	"net/http"

	"github.com/ddollar/logger"
	"github.com/gorilla/mux"
)

func New(ns, hostname string) *Server {
	logger := logger.New(fmt.Sprintf("ns=%s", ns))

	server := &Server{
		Hostname: hostname,
		Logger:   logger,
	}

	server.Router = &Router{
		Parent: nil,
		Router: mux.NewRouter(),
		Server: server,
	}

	server.Router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		id, _ := generateId(12)
		logger.Logf("id=%s route=unknown code=404 method=%q path=%q", id, r.Method, r.URL.Path)
	})

	server.Router.HandleFunc("/check", server.check)

	return server
}

func (s *Server) check(w http.ResponseWriter, r *http.Request) {
	if s.Check != nil {
		if err := s.Check(NewContext(w, r)); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}

	fmt.Fprintf(w, "ok")
}
