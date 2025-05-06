package limiter

import (
	"RateBalancer/internal/handler/http/api"
	"RateBalancer/internal/handler/http/middleware"
	"RateBalancer/internal/service"
	"log/slog"
	"net/http"
)

type Limiter struct {
	log            *slog.Logger
	limiterService service.Limiter
}

func NewLimiter(ls service.Limiter, log *slog.Logger) *Limiter {
	return &Limiter{
		log:            log,
		limiterService: ls,
	}
}

func (l *Limiter) RegisterLimiter(h http.Handler) http.Handler {
	return l.checkLimit(h)
}

func (l *Limiter) checkLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			l.log.Error(NoApiKey.Error())
			api.ErrorResponseWithCode(w, r, http.StatusUnauthorized, NoApiKey)
			return
		}

		allowed, err := l.limiterService.ConsumeTokens(r.Context(), apiKey)
		if err != nil {
			l.log.Error("failed to check limit", slog.String("error", err.Error()))
			api.ErrorResponse(w, r, err)
			return
		}

		if allowed {
			l.log.Info("Consume token successful")
			next.ServeHTTP(w, r)
		} else {
			l.log.Error(RateLimitExceeded.Error(),
				slog.String("request", r.URL.String()),
				slog.String("request_id", middleware.GetRequestIdFromContext(r)))
			api.ErrorResponseWithCode(w, r, http.StatusTooManyRequests, RateLimitExceeded)
		}
	})
}
