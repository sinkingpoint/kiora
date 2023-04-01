package frontend

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gorilla/mux"
)

//go:generate npm --prefix ../../../frontend run build
//go:generate rm -r assets
//go:generate mkdir -p assets
//go:generate cp -r ../../../frontend/build assets
//go:embed assets/build/*
var content embed.FS

// Registers a router that will serve the frontend assets created and embedded by the above go:generate directives.
func Register(router *mux.Router) {
	sub, err := fs.Sub(content, "assets/build")

	if err != nil {
		panic("BUG: failed to embed frontend assets: " + err.Error())
	}

	router.PathPrefix("/").Handler(http.FileServer(http.FS(sub)))
}