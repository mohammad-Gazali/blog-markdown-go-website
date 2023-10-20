package router

import (
	"net/http"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	serveStaticFiles(mux, "static")

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		MustGetTemplate(w, "index.html", nil)
	})

	return mux
}

func serveStaticFiles(mux *http.ServeMux, folderName string) {
	fs := http.FileServer(http.Dir(folderName))
	mux.Handle("/" + folderName + "/", http.StripPrefix("/" + folderName, fs))
}