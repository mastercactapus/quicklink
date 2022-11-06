package web

import (
	_ "embed"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/mastercactapus/quicklink/pkg/store"
)

//go:embed index.html
var page string

var tmpl = template.Must(template.New("index").Parse(page))

type App struct {
	Store store.Store
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		a.handlePost(w, r)
		return
	}

	type link struct {
		Base string
		Dest string
	}
	s, err := a.Store.Scanner(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer s.Close()

	var links []link
	for s.Next() {
		links = append(links, link{
			Base: s.Key(),
			Dest: s.Value(),
		})
	}
	err = s.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data struct {
		Links []link
	}
	data.Links = links

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) handlePost(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/_/delete":
		err := a.Store.Set(r.Context(), r.FormValue("base"), "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
		return
	case "/_/add":
		base := r.FormValue("base")
		dest := r.FormValue("dest")
		if base == "" || dest == "" {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}
		if strings.Contains(base, "/") {
			http.Error(w, "invalid base", http.StatusBadRequest)
			return
		}
		u, err := url.Parse(dest)
		if err != nil || u.Scheme == "" {
			http.Error(w, "invalid destination", http.StatusBadRequest)
			return
		}

		err = a.Store.Set(r.Context(), base, dest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/#"+base, http.StatusFound)
		return
	}

	http.NotFound(w, r)
}
