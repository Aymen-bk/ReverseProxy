# Concurrent Load-Balancing Reverse Proxy

A high-performance reverse proxy server built in Go that implements load balancing, health monitoring, and dynamic backend management.

## Features

- **Round-Robin Load Balancing**: Distributes incoming requests evenly across healthy backends
- **Health Monitoring**: Periodic background health checks to ensure backend availability
- **Thread-Safe Operations**: Uses `sync.RWMutex` and `sync/atomic` for concurrent safety
- **Context Propagation**: Proper request context handling for cancellation and timeouts
- **Admin API**: RESTful API for managing backends dynamically
- **Graceful Shutdown**: Handles shutdown signals gracefully

## Project Structure

```
Reverse_Proxy_Aymen/
├── main.go              # Main entry point
├── Backend.go           # Backend data structure and health checking
├── ServerPool.go        # Server pool with Round-Robin load balancing
├── LoadBalancer.go      # LoadBalancer interface
├── ProxyHandler.go      # HTTP request handler and proxy logic
├── ProxyConfig.go       # Configuration loading
├── HealthChecker.go     # Background health monitoring
├── AdminApi.go          # Admin API endpoints
├── config.json          # Configuration file
├── go.mod               # Go module file
└── README.md           # This file
```

## Installation

1. Make sure you have Go 1.21 or later installed
2. Clone or navigate to the project directory
3. Initialize dependencies:
   ```bash
   go mod tidy
   ```

## Configuration

Edit `config.json` to configure the proxy:

```json
{
  "port": 8000,
  "admin_port": 8081,
  "strategy": "round-robin",
  "health_check_frequency": "30s",
  "backends": [
    "http://localhost:8001",
    "http://localhost:8002",
    "http://localhost:8003"
  ]
}
```

### Configuration Options

- `port`: Port for the proxy server (default: 8000)
- `admin_port`: Port for the admin API (default: 8081)
- `strategy`: Load balancing strategy (currently supports "round-robin")
- `health_check_frequency`: How often to check backend health (e.g., "30s", "1m")
- `backends`: Initial list of backend URLs

## Usage

### Starting the Proxy

```bash
go run main.go --config=config.json
```

Or build and run:

```bash
go build -o reverse_proxy
./reverse_proxy --config=config.json
```

### Example Backend Servers

You can create simple test backends. Example (`backend1.go`):

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hi from server 1")
    })
    fmt.Println("Server is running on http://localhost:8001")
    http.ListenAndServe(":8001", nil)
}
```

Run multiple backends on different ports (8001, 8002, 8003, etc.)

## Admin API

The admin API runs on a separate port (default: 8081) and provides endpoints for managing backends.

### GET /status

Get the status of all backends:

```bash
curl http://localhost:8081/status
```

Response:
```json
{
  "total_backends": 3,
  "active_backends": 2,
  "backends": [
    {
      "url": "http://localhost:8001",
      "alive": true,
      "current_connections": 5
    },
    {
      "url": "http://localhost:8002",
      "alive": true,
      "current_connections": 3
    },
    {
      "url": "http://localhost:8003",
      "alive": false,
      "current_connections": 0
    }
  ]
}
```

### POST /backends

Add a new backend to the pool:

```bash
curl -X POST http://localhost:8081/backends \
  -H "Content-Type: application/json" \
  -d '{"url": "http://localhost:8004"}'
```

### DELETE /backends

Remove a backend from the pool:

```bash
curl -X DELETE http://localhost:8081/backends \
  -H "Content-Type: application/json" \
  -d '{"url": "http://localhost:8004"}'
```

## Testing

1. Start multiple backend servers on different ports
2. Start the proxy server
3. Send requests to the proxy (port 8000)
4. Check the admin API to see load distribution
5. Stop a backend server and observe health checker detecting it

Example test:

```bash
# Terminal 1: Start backend 1
go run backend1.go

# Terminal 2: Start backend 2
go run backend2.go

# Terminal 3: Start the proxy
go run main.go --config=config.json

# Terminal 4: Send requests
curl http://localhost:8000/
curl http://localhost:8000/
curl http://localhost:8000/

# Check status
curl http://localhost:8081/status
```

## Architecture

### Components

1. **Backend**: Represents a backend server with URL, alive status, and connection count
2. **ServerPool**: Manages a collection of backends with Round-Robin selection
3. **ProxyHandler**: Handles incoming HTTP requests and forwards them to backends
4. **HealthChecker**: Periodically checks backend health in the background
5. **AdminApi**: Provides REST endpoints for backend management

### Concurrency Safety

- `sync.RWMutex` is used for read-write locks on shared data structures
- `sync/atomic` is used for atomic operations on counters
- Each backend has its own mutex for thread-safe status updates
- ServerPool uses mutexes to protect backend list modifications

### Load Balancing

The proxy implements Round-Robin load balancing:
- Uses an atomic counter to track the current backend index
- Skips dead backends automatically
- Returns 503 Service Unavailable if no backends are available

### Health Checking

- Runs in a background goroutine
- Configurable check interval (default: 30 seconds)
- Performs HTTP GET requests to verify backend availability
- Updates backend status automatically
- Logs status changes

## Error Handling

- Backend connection failures are caught and the backend is marked as dead
- Client disconnections are handled via context cancellation
- Invalid backend URLs are logged and skipped
- Graceful shutdown handles in-flight requests

## Future Enhancements

- Least-Connections load balancing strategy
- Weighted load balancing
- Sticky sessions (session affinity)
- HTTPS/TLS support
- Metrics and monitoring endpoints
- Request/response logging

## License

This project is part of a Go programming course final project.

