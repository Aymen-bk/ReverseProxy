package main

import (
	"log"
	"time"
)

// HealthChecker performs periodic health checks on all backends
type HealthChecker struct {
	serverPool *ServerPool
}

// Start begins the health checking process in a background goroutine
func (hc *HealthChecker) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			log.Println("Starting health check...")
			hc.CheckAllBackends()
			log.Println("Health check completed")
		}
	}()
}

// CheckAllBackends checks the health status of all backends in the pool
func (hc *HealthChecker) CheckAllBackends() {
	hc.serverPool.mux.RLock()
	backends := make([]*Backend, len(hc.serverPool.Backends))
	copy(backends, hc.serverPool.Backends)
	hc.serverPool.mux.RUnlock()

	for _, backend := range backends {
		status := "UP"
		alive := backend.GetRealStatus()
		backend.SetAlive(alive)
		
		if !alive {
			status = "DOWN"
		}
		log.Printf("%s is %s (connections: %d)", backend.URL.String(), status, backend.CurrentConns.Load())
	}
}

