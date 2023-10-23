package main

import (
	"blog-markdown-website/server"
	"flag"
	"log"
	"net"
)

func main() {
	var (
		host = flag.String("host", "", "host http address to listen on")
		port = flag.String("port", "8000", "port number for http listener")
	)
	flag.Parse()

	addr := net.JoinHostPort(*host, *port)

	if err := server.CreateAndRunServer(addr); err != nil {
		log.Fatal(err)
	}
}