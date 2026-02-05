package main

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	URL          *url.URL `json:"url"`
	Alive        bool     `json:"alive"`
	CurrentConns int64    `json:"current_connections"`
	mux          sync.RWMutex
}

func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Alive = alive
}

func (b *Backend) GetRealStatus() bool {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(b.URL.String())
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode < 500
}
