package main
a
import (
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"
)

type ServerPool struct {
	Backends []*Backend `json:"backends"`
	Current  uint64     `json:"current"` 
	mux      sync.RWMutex
}

func (sp *ServerPool) GetNextValidPeer() *Backend {
	sp.mux.RLock()
	defer sp.mux.RUnlock()

	if len(sp.Backends)==0 { //skip
		return nil
	}

	rrindex :=atomic.AddUint64(&sp.Current,1) % uint64(len(sp.Backends))

	for i := 0; i < len(sp.Backends); i++ {
		idx:= (rrindex+uint64(i))% uint64(len(sp.Backends))
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

func (sp *ServerPool) AddBackend(backend *Backend) {
	sp.mux.Lock()
	sp.Backends = append(sp.Backends, backend)
}

func (sp *ServerPool) SetBackendStatus(uri *url.URL, alive bool) {
	sp.mux.Lock()
	defer sp.mux.Unlock()
	for _ , backend := range sp.Backends {
		if backend.URL.String() == uri.String() {
			backend.mux.Lock()
			backend.Alive = alive
			backend.mux.Unlock()
			return
		}
	}
}
