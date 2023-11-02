package server

import (
	"fmt"
	"net/http"
)


type RenderContext map[string]any


func CreateAndRunServer(listenAddr string) error {
	fmt.Printf("Starting HTTP Server at %v\n", listenAddr)
	ServeStaticFiles("static")
	AddingHandlers()
	return http.ListenAndServe(listenAddr, nil)
}


// ======== utils ========

func ServeStaticFiles(folderName string) {
	fs := http.FileServer(http.Dir(folderName))
	
	http.Handle("/" + folderName + "/", http.StripPrefix("/" + folderName, fs))
}

func InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	RenderTemplate(w, RenderContext{"Title": "Server Error"}, "500.html")
}

func NotFoundError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	RenderTemplate(w, RenderContext{"Title": "Not Found"}, "404.html")
}