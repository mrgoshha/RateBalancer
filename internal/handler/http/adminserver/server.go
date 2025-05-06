package adminserver

import (
	"RateBalancer/internal/handler/http/middleware"
	"fmt"
	"log/slog"
	"net/http"
)

func NewServer(cfg *Config, logger *slog.Logger, handler http.Handler) *http.Server {
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	srv := &http.Server{
		Addr:    address,
		Handler: middleware.SetRequestId(middleware.NewLogger(logger).LogRequest(handler)),
	}

	return srv
}
