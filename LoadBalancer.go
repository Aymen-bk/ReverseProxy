package main

import (
	"net/url"
)

// LoadBalancer interface defines methods for load balancing operations
type LoadBalancer interface {
	GetNextValidPeer() *Backend
	AddBackend(backend *Backend)
	SetBackendStatus(uri *url.URL, alive bool)
	RemoveBackend(uri *url.URL)
}

