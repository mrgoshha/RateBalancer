package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type Logger struct {
	log *slog.Logger
}

func NewLogger(log *slog.Logger) *Logger {
	return &Logger{
		log: log,
	}
}

func (l *Logger) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := l.log.With(
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
			slog.String("request_id", GetRequestIdFromContext(r)),
		)

		logger.Info("request started")

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		logger.Info("request completed with",
			slog.Int("status", rw.code),
			slog.String("duration", time.Since(start).String()),
		)
	})
}
