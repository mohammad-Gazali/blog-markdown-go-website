package server

import "net/http"

func AddingHandlers(mux *http.ServeMux) {
	HandleFuncWith404(mux, "/", HomeHandler)
}

// handlers
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, RenderContext{"Title": "Markdown Blog"}, "index.html")
}