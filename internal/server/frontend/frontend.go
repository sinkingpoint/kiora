package frontend

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

//go:embed assets/*
var content embed.FS

// Registers a router that will serve the frontend assets created and embedded by the above go:generate directives.
func Register(router *mux.Router) {
	sub, err := fs.Sub(content, "assets")
	if err != nil {
		panic("BUG: failed to embed frontend assets: " + err.Error())
	}

	serveReactApp := func(w http.ResponseWriter, r *http.Request) {
		idx, err := fs.ReadFile(sub, "index.html")
		if err != nil {
			panic("BUG: failed to read embedded index.html: " + err.Error())
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write(idx)
	}

	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, ".") {
			http.FileServer(http.FS(sub)).ServeHTTP(w, r)
		} else {
			serveReactApp(w, r)
		}
	})
}
