package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func main() {
	configFile := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	config, err := LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	serverPool := &ServerPool{
		Backends: make([]*Backend, 0),
		Current:  0,
	}

	for _, backendURL := range config.Backends {
		u, err := url.Parse(backendURL)
		if err != nil {
			log.Printf("Warning: Invalid backend URL %s: %v", backendURL, err)
			continue
		}
		backend := &Backend{
			URL:   u,
			Alive: true,
		}
		serverPool.AddBackend(backend)
		log.Printf("Added backend: %s", u.String())
	}

	handler := &ProxyHandler{
		lb: serverPool,
	}

	healthChecker := &HealthChecker{
		serverPool: serverPool,
	}
	healthChecker.Start(config.HealthCheckFreq)

	go func() {
		StartAdminServer(serverPool, config.AdminPort)
	}()

	fmt.Printf("Proxy server listening on :%d\n", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), handler))
}
