package server

import (
	"net/http"
)

type RenderContext map[string]any

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	ServeStaticFiles(mux, "static")

	AddingHandlers(mux)

	return mux
}

func ServeStaticFiles(mux *http.ServeMux, folderName string) {
	fs := http.FileServer(http.Dir(folderName))
	
	mux.Handle("/" + folderName + "/", http.StripPrefix("/" + folderName, fs))
}
