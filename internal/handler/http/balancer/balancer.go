package balancer

import (
	"RateBalancer/internal/handler/http/api"
	"RateBalancer/internal/handler/http/middleware"
	"RateBalancer/internal/service"
	bs "RateBalancer/internal/service/balancer"
	"log/slog"
	"net/http"
)

type Balancer struct {
	log         *slog.Logger
	strategy    service.Strategy
	backendPool *bs.BackendPool
}

func NewBalancer(strategy service.Strategy, backendPool *bs.BackendPool, log *slog.Logger) *Balancer {
	balancer := &Balancer{
		log:         log,
		strategy:    strategy,
		backendPool: backendPool,
	}

	for _, b := range backendPool.Backends {
		b.ReverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
			log.Error("failed to send request",
				slog.String("request", r.URL.String()),
				slog.String("request_id", middleware.GetRequestIdFromContext(r)),
				slog.String("server", b.URL.String()),
				slog.String("error", e.Error()))

			b.HandleFailure()
			balancer.balance(w, r)
		}
	}

	return balancer
}

func (b *Balancer) RegisterBalancer(router *http.ServeMux) *http.ServeMux {
	router.HandleFunc("/", b.balance)
	return router
}

func (b *Balancer) balance(w http.ResponseWriter, r *http.Request) {
	backend, err := b.strategy.GetNext()
	if err != nil {
		b.log.Error(ServiceNotAvailable.Error(), slog.String("error", err.Error()))
		api.ErrorResponseWithCode(w, r, http.StatusServiceUnavailable, ServiceNotAvailable)
		return
	}

	if backend == nil {
		b.log.Error(ServerNotExist.Error())
		api.ErrorResponseWithCode(w, r, http.StatusInternalServerError, ServerNotExist)
		return
	}

	b.log.Info("balance request to server",
		slog.String("request", r.URL.String()),
		slog.String("request_id", middleware.GetRequestIdFromContext(r)),
		slog.String("server", backend.URL.String()))
	backend.ReverseProxy.ServeHTTP(w, r)
}
