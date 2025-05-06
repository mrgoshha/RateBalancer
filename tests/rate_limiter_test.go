package tests

import (
	"RateBalancer/internal/handler/http/limiter"
	"RateBalancer/internal/handler/http/middleware"
	"RateBalancer/internal/handler/http/model"
	servicelimiter "RateBalancer/internal/service/limiter"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"time"
)

func (s *TestSuite) TestLimiter() {
	//Arrange
	router := http.NewServeMux()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Handle("/", handler)
	limitedHandler := s.limiter.RegisterLimiter(router)

	// Create Request
	req := func() *http.Request {
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("X-API-Key", "a2d28bc8-b419-4e51-837a-fe50bac1c991")
		ctx := context.WithValue(r.Context(), middleware.CtxKeyRequestID, "req-4")
		return r.WithContext(ctx)
	}
	// capacity = 2 per_second = 1 token = 0

	query := ` UPDATE clients
        	  SET last_refill = $1
        	  WHERE api_key = $2`
	_, err := s.db.Exec(query, time.Now().UTC(), "0c2975ce3a7af8aef09a40655c38129822c5d074")
	s.Require().NoError(err, "failed to update last_refill")

	time.Sleep(time.Second * 2)

	// Make Request and Assert
	for i := 0; i < 2; i++ {
		rec := httptest.NewRecorder()
		limitedHandler.ServeHTTP(rec, req())
		s.Equal(http.StatusOK, rec.Code)
	}

	// Запрос должен превысить лимит → 429
	rec := httptest.NewRecorder()
	limitedHandler.ServeHTTP(rec, req())
	s.Equal(http.StatusTooManyRequests, rec.Code)

	expectedBody := &model.ErrorResponse{
		Code:    429,
		Message: limiter.RateLimitExceeded.Error(),
	}
	var body *model.ErrorResponse
	err = json.NewDecoder(rec.Body).Decode(&body)
	s.Require().NoError(err, "decode error response")
	s.Equal(expectedBody.Message, body.Message)

	time.Sleep(time.Second + 50*time.Millisecond)

	rec = httptest.NewRecorder()
	limitedHandler.ServeHTTP(rec, req())
	s.Equal(http.StatusOK, rec.Code)
}

func (s *TestSuite) TestNotFoundUser() {
	//Arrange
	router := http.NewServeMux()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Handle("/", handler)
	limitedHandler := s.limiter.RegisterLimiter(router)

	// Create Request
	req := func() *http.Request {
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("X-API-Key", "111")
		ctx := context.WithValue(r.Context(), middleware.CtxKeyRequestID, "req-5")
		return r.WithContext(ctx)
	}

	rec := httptest.NewRecorder()
	limitedHandler.ServeHTTP(rec, req())
	s.Equal(http.StatusUnauthorized, rec.Code)

	expectedBody := &model.ErrorResponse{
		Code:    401,
		Message: servicelimiter.InvalidAPIKey.Error(),
	}
	var body *model.ErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&body)
	s.Require().NoError(err, "decode error response")
	s.Contains(body.Message, expectedBody.Message)
}
