package balancer

import (
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

type BackendPool struct {
	Backends []*Backend
	Current  atomic.Uint64
}

func NewBackendPool(cfg *Config) *BackendPool {
	pool := &BackendPool{
		Backends: make([]*Backend, len(cfg.Backends)),
	}
	pool.Current.Store(0)

	for i, b := range cfg.Backends {
		proxy := httputil.NewSingleHostReverseProxy(b.URL)

		backend := &Backend{
			URL:                b.URL,
			ReverseProxy:       proxy,
			HealthyThreshold:   uint32(cfg.HealthyThreshold),
			UnhealthyThreshold: uint32(cfg.UnhealthyThreshold),
			Timeout:            cfg.PingTimeout,
		}

		backend.SetAlive(true)

		pool.Backends[i] = backend
	}

	return pool
}

func (s *BackendPool) Ping() {
	wg := sync.WaitGroup{}
	for _, b := range s.Backends {
		if !b.Alive() {
			wg.Add(1)

			go func(b *Backend, url *url.URL) {
				defer wg.Done()
				if b.IsAlive(b.URL) {
					b.HandleSuccess()
				}
			}(b, b.URL)
		}

	}

	wg.Wait()
}
