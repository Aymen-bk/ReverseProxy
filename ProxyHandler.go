package main

import (
	"net/http"
	"net/http/httputil"
)

// ProxyHandler handles incoming HTTP requests and forwards them to backends
type ProxyHandler struct {
	lb LoadBalancer
}

// ServeHTTP implements the http.Handler interface
func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the next available backend using the load balancer
	backend := p.lb.GetNextValidPeer()
	if backend == nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	// Increase active connections atomically
	backend.CurrentConns.Add(1)
	defer backend.CurrentConns.Add(-1)

	// Create a reverse proxy for this backend
	proxy := httputil.NewSingleHostReverseProxy(backend.URL)

	// Preserve the request context so cancellations propagate
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req = req.WithContext(r.Context())
	}

	// Handle backend errors - mark backend as dead if connection fails
	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		backend.SetAlive(false)
		http.Error(w, "Backend unavailable", http.StatusBadGateway)
	}

	// Forward the request to the backend
	proxy.ServeHTTP(w, r)
}

