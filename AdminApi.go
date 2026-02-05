package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func StartAdminServer(serverPool *ServerPool, port int) {
	handleGetStatus := func(w http.ResponseWriter, r *http.Request) {
		serverPool.mux.RLock()
		backends := make([]*Backend, len(serverPool.Backends))
		copy(backends, serverPool.Backends)
		serverPool.mux.RUnlock()

		activeCount := 0
		for _, b := range backends {
			b.mux.RLock()
			if b.Alive {
				activeCount++
			}
			b.mux.RUnlock()
		}

		response := map[string]interface{}{
			"total_backends":  len(backends),
			"active_backends": activeCount,
			"backends":        backends,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}

	handlePostBackends := func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u, err := url.Parse(req.URL)
		if err != nil {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}
		backend := &Backend{
			URL:   u,
			Alive: true,
		}
		serverPool.AddBackend(backend)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Backend added successfully"})
	}

	handleDeleteBackends := func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u, err := url.Parse(req.URL)
		if err != nil {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}
		serverPool.RemoveBackend(u)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Backend removed successfully"})
	}

	http.HandleFunc("/status", handleGetStatus)
	http.HandleFunc("/backends", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlePostBackends(w, r)
		} else if r.Method == http.MethodDelete {
			handleDeleteBackends(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Printf("Admin server listening on :%d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("Failed to start admin server: %v", err)
	}
}
