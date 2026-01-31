package main

import (
	"net/http"
	"net/http/httputil"
)

type ProxyHandler struct {
	lb LoadBalancer
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := p.lb.GetNextValidPeer()
	if backend == nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	backend.AddConns(1)
	defer backend.AddConns(-1)

	proxy := httputil.NewSingleHostReverseProxy(backend.URL)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req = req.WithContext(r.Context())
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		backend.SetAlive(false)
		http.Error(w, "Backend unavailable", http.StatusBadGateway)
	}

	proxy.ServeHTTP(w, r)
}
