package server

import (
	"fmt"
	"net/http"
)

func CreateAndRunServer(listenAddr string) error {
	s := http.Server{
		Addr:    listenAddr,
		Handler: NewRouter(),
	}

	fmt.Printf("Starting HTTP Server at %v\n", listenAddr)

	return s.ListenAndServe()
}