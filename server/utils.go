package server

import "net/http"

// this function wrap the handler with 404 Error and Page
func HandleFuncWith404(mux *http.ServeMux, route string, hanlder func(w http.ResponseWriter, r *http.Request)) {
	mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		if route != r.URL.Path {
			w.WriteHeader(http.StatusNotFound)
			RenderTemplate(w, RenderContext{"Title": "Not Found"}, "404.html")
			return
		}
		hanlder(w, r)
	})
}

func InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	RenderTemplate(w, RenderContext{"Title": "Server Error"}, "500.html")
}
