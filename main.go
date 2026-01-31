package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const listenPort = ":9000"

func mustParseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		log.Fatalf("invalid backend URL %q: %v", raw, err)
	}
	return u
}

func newPoolFromURLs(rawURLs []string) *ServerPool {
	pool := &ServerPool{Current: 0}
	for _, raw := range rawURLs {
		pool.Backends = append(pool.Backends, &Backend{
			URL:   mustParseURL(raw),
			Alive: true,
		})
	}
	return pool
}

func main() {
	upstreamURLs := []string{
		"http://localhost:8001",
		"http://localhost:8002",
	}
	pool := newPoolFromURLs(upstreamURLs)
	proxy := &ProxyHandler{lb: pool}

	fmt.Printf("Proxy listening on %s\n", listenPort)
	log.Fatal(http.ListenAndServe(listenPort, proxy))
}
