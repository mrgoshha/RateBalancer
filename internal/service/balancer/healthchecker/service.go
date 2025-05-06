package healthchecker

import (
	"RateBalancer/internal/service/balancer"
	"context"
	"time"
)

type HealthChecker struct {
	bp           *balancer.BackendPool
	pingInterval time.Duration
}

func NewHealthChecker(bp *balancer.BackendPool, cfg *Config) *HealthChecker {
	return &HealthChecker{
		bp:           bp,
		pingInterval: cfg.PingInterval,
	}
}

func (h *HealthChecker) HealthCheck(ctx context.Context) {
	ticker := time.NewTicker(h.pingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			h.bp.Ping()
		case <-ctx.Done():
			return
		}
	}
}
