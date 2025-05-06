package balancer

import (
	"net"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

type Backend struct {
	URL                *url.URL
	alive              atomic.Bool
	ReverseProxy       *httputil.ReverseProxy
	UnhealthyThreshold uint32
	failCount          atomic.Uint32
	HealthyThreshold   uint32
	successCount       atomic.Uint32
	Timeout            time.Duration
}

func (b *Backend) SetAlive(alive bool) {
	b.alive.Store(alive)
}

func (b *Backend) Alive() bool {
	return b.alive.Load()
}

func (b *Backend) IsAlive(u *url.URL) bool {
	conn, err := net.DialTimeout("tcp", u.Host, b.Timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func (b *Backend) HandleFailure() {
	if !b.Alive() {
		return
	}
	b.successCount.Store(0)
	newFails := b.failCount.Add(1)

	if b.Alive() && newFails >= b.UnhealthyThreshold {
		b.alive.CompareAndSwap(true, false)
	}
}

func (b *Backend) HandleSuccess() {
	if b.Alive() {
		return
	}
	b.failCount.Store(0)
	newSuccess := b.successCount.Add(1)

	if !b.Alive() && newSuccess >= b.HealthyThreshold {
		b.alive.CompareAndSwap(false, true)
	}
}
