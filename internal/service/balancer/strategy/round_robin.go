package strategy

import (
	"RateBalancer/internal/service/balancer"
)

type RoundRobin struct {
	bp *balancer.BackendPool
}

func NewRoundRobinBalancer(bp *balancer.BackendPool) *RoundRobin {
	return &RoundRobin{
		bp: bp,
	}
}

func (r *RoundRobin) GetNext() (*balancer.Backend, error) {
	lenBackends := len(r.bp.Backends)
	next := int(r.bp.Current.Add(1) % uint64(lenBackends))

	for offset := 0; offset < lenBackends; offset++ {
		idx := (next + offset) % lenBackends
		if r.bp.Backends[idx].Alive() {
			if offset != 0 {
				r.bp.Current.Store(uint64(idx))
			}
			return r.bp.Backends[idx], nil
		}
	}

	return nil, AllServiceNotAvailable
}
