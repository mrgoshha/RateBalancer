package strategy

import (
	"RateBalancer/internal/service/balancer"
	"math/rand"
	"time"
)

type Random struct {
	bp  *balancer.BackendPool
	rnd *rand.Rand
}

func NewRandomBalancer(bp *balancer.BackendPool) *Random {
	src := rand.NewSource(time.Now().UnixNano())
	return &Random{
		bp:  bp,
		rnd: rand.New(src),
	}
}

func (r *Random) GetNext() (*balancer.Backend, error) {
	lenBackends := len(r.bp.Backends)
	next := r.rnd.Intn(lenBackends)

	for offset := 0; offset < lenBackends; offset++ {
		idx := (next + offset) % lenBackends
		if r.bp.Backends[idx].Alive() {
			return r.bp.Backends[idx], nil
		}
	}

	return nil, AllServiceNotAvailable
}
