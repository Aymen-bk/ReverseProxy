package main

import (
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"
)

// ServerPool implements the LoadBalancer interface using Round-Robin strategy
type ServerPool struct {
	Backends []*Backend `json:"backends"`
	Current  uint64     `json:"current"` // Used for Round-Robin
	mux      sync.RWMutex
}

// GetNextValidPeer returns the next available backend using Round-Robin algorithm
// It skips dead backends and returns nil if no backends are available
func (sp *ServerPool) GetNextValidPeer() *Backend {
	sp.mux.RLock()
	defer sp.mux.RUnlock() // safe read lock for backends

	if len(sp.Backends) == 0 {
		return nil
	}

	// Use atomic operation for thread-safe counter increment
	rrindex := atomic.AddUint64(&sp.Current, 1) % uint64(len(sp.Backends))

	// Try to find an alive backend starting from the round-robin index
	for i := 0; i < len(sp.Backends); i++ {
		idx := (rrindex + uint64(i)) % uint64(len(sp.Backends))
		backend := sp.Backends[idx]

		backend.mux.RLock()
		alive := backend.Alive
		backend.mux.RUnlock()

		if alive {
			return backend
		}
	}

	fmt.Println("No Server Available")
	return nil
}

// AddBackend adds a new backend to the server pool
func (sp *ServerPool) AddBackend(backend *Backend) {
	sp.mux.Lock()
	defer sp.mux.Unlock()
	sp.Backends = append(sp.Backends, backend)
}

// SetBackendStatus updates the alive status of a backend by its URL
func (sp *ServerPool) SetBackendStatus(uri *url.URL, alive bool) {
	sp.mux.Lock()
	defer sp.mux.Unlock()
	
	for _, backend := range sp.Backends {
		if backend.URL.String() == uri.String() {
			backend.SetAlive(alive)
			return
		}
	}
}

// RemoveBackend removes a backend from the server pool by its URL
func (sp *ServerPool) RemoveBackend(backendURL *url.URL) {
	sp.mux.Lock()
	defer sp.mux.Unlock()
	
	for i, v := range sp.Backends {
		if v.URL.String() == backendURL.String() {
			sp.Backends = append(sp.Backends[:i], sp.Backends[i+1:]...)
			return
		}
	}
}

