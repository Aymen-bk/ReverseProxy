package main

import (
	"fmt"
	"log"
	"net"
)
func parseURL(raw string) *net.URL {
	urll, err := url.Parse(raw)
	if err != nil {
		log.Fatal(err)
	}
	return urll
}

func main() {
	backend1:=&Backend{URL: parseURL("http://localhost:9001"), Alive: true}
	backend2:=&Backend{URL: parseURL("http://localhost:9002"), Alive: true}
	serverPool := &ServerPool{
		Backends: []*Backend{backend1, backend2},
		Current:  0,
	}
	handler := &ProxyHandler{lb: serverPool}
	fmt.Println("Running on :9000")
	log.Fatal(http.ListenAndServe(":9000", handler))
}
