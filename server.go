package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/mastercactapus/quicklink/internal/web"
	"github.com/mastercactapus/quicklink/pkg/store"
)

type Server struct {
	store.Store

	app *web.App
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	if r.URL.Path == "/" || strings.HasPrefix(r.URL.Path, "/_/") {
		s.app.ServeHTTP(w, r)
		return
	}

	id, _, _ := strings.Cut(r.URL.Path[1:], "/")
	val, err := s.Store.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if val == "" {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, val, http.StatusFound)
}
