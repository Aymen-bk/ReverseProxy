package main

import (
	"log"
	"time"
)

type HealthChecker struct {
	serverPool *ServerPool
}

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
		backend.mux.RLock()
		conns := backend.CurrentConns
		backend.mux.RUnlock()
		log.Printf("%s is %s (connections: %d)", backend.URL.String(), status, conns)
	}
}
