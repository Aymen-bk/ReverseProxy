package main

import (
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type Backend struct {
	URL          *url.URL    `json:"url"`
	Alive        bool        `json:"alive"`
	CurrentConns atomic.Int64 `json:"current_connections"`
	mux          sync.RWMutex
}

// SetAlive safely updates the alive status of the backend
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Alive = alive
}

// GetRealStatus performs an actual HTTP request to check if the backend is alive
func (b *Backend) GetRealStatus() bool {
	client := &http.Client{
		Timeout: 2 * time.Second, // wait for 2 seconds
	}

	resp, err := client.Get(b.URL.String())
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode < 500 // Consider 2xx, 3xx, 4xx as "alive"
}

