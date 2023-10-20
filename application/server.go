package application

import (
	"blog-markdown-website/router"
	"fmt"
	"net/http"
)

func RunServer(listenAddr string) error {
	s := http.Server{
		Addr:    listenAddr,
		Handler: router.NewRouter(),
	}

	fmt.Printf("Starting HTTP Server at %v\n", listenAddr)

	return s.ListenAndServe()
}