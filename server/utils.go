package server

import (
	"net/http"
	"strings"
)

// ============================ Web Utils ============================

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


// ============================ Markdown Utils ============================

// this function return the info in string line
//
// for example:
//
// - GetLineInfo("name: Mohammed Algazali", "name")    >> "Mohammed Algazali"
//
// - GetLineInfo("created_at: 2023/4/1", "created_at") >> "2023/4/1"
//
// - GetLineInfo("job: programmer", "jop")             >> ""  # because we failed in key name writing
// 
func GetLineInfo(line, key string) string {
	if info := strings.Split(line, key); len(info) >= 2 {
		// here we remove ":", spaces and "\r" characters
		return strings.Trim(info[1], " " + ":" + "\r")
	} else {
		return ""
	}
}