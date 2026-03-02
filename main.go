package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Parse command line arguments
	configFile := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	// Load configuration
	config, err := LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize server pool
	serverPool := &ServerPool{
		Backends: make([]*Backend, 0),
		Current:  0,
	}

	// Add initial backends from config
	for _, backendURL := range config.Backends {
		u, err := url.Parse(backendURL)
		if err != nil {
			log.Printf("Warning: Invalid backend URL %s: %v", backendURL, err)
			continue
		}

		backend := &Backend{
			URL:   u,
			Alive: true, // Assume alive initially, health checker will verify
		}
		serverPool.AddBackend(backend)
		log.Printf("Added backend: %s", u.String())
	}

	// Create proxy handler
	handler := &ProxyHandler{
		lb: serverPool,
	}

	// Initialize health checker
	healthChecker := &HealthChecker{
		serverPool: serverPool,
	}

	// Start health checking in background
	healthChecker.Start(config.HealthCheckFreq)

	// Start admin API server in background
	go func() {
		StartAdminServer(serverPool, config.AdminPort)
	}()

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start proxy server in a goroutine
	go func() {
		log.Printf("Proxy server listening on :%d", config.Port)
		log.Printf("Strategy: %s", config.Strategy)
		log.Printf("Health check frequency: %v", config.HealthCheckFreq)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start proxy server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	log.Println("Shutting down gracefully...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("Server stopped")
}

