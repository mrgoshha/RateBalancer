package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

const (
	CtxKeyRequestID = "requestId"
)

func SetRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKeyRequestID, id)))
	})
}

func GetRequestIdFromContext(r *http.Request) string {
	if requestId, ok := r.Context().Value(CtxKeyRequestID).(string); ok {
		return requestId
	}
	return ""
}
