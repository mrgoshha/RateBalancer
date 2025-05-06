package tests

import (
	"RateBalancer/internal/handler/http/balancer"
	"RateBalancer/internal/handler/http/middleware"
	"RateBalancer/internal/handler/http/model"
	"RateBalancer/internal/service/balancer/strategy"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func (s *TestSuite) TestAllAlive() {
	//Arrange
	router := http.NewServeMux()
	backendPool := s.newBackendPool(s.urls)
	strategyRR := strategy.NewRoundRobinBalancer(backendPool)
	b := balancer.NewBalancer(strategyRR, backendPool, s.log)
	b.RegisterBalancer(router)

	for idx, want := range []string{"Server1", "Server2", "Server0", "Server1"} {
		// Create Request
		ctx := context.WithValue(context.Background(), middleware.CtxKeyRequestID, "req-1")
		req := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
		rec := httptest.NewRecorder()

		// Make Request
		router.ServeHTTP(rec, req)

		//Assert
		s.Equal(http.StatusOK, rec.Code, "step %d: expected 200", idx)
		s.Equal(want, rec.Body.String(), "step %d: body", idx)
	}

}

func (s *TestSuite) TestSkipDeadServer() {
	//Arrange
	router := http.NewServeMux()
	backendPool := s.newBackendPool(s.urls)
	strategyRR := strategy.NewRoundRobinBalancer(backendPool)
	b := balancer.NewBalancer(strategyRR, backendPool, s.log)
	b.RegisterBalancer(router)

	backendPool.Backends[1].SetAlive(false)

	for idx, want := range []string{"Server2", "Server0", "Server2", "Server0"} {
		// Create Request
		ctx := context.WithValue(context.Background(), middleware.CtxKeyRequestID, "req-2")
		req := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
		rec := httptest.NewRecorder()

		// Make Request
		router.ServeHTTP(rec, req)

		//Assert
		s.Equal(http.StatusOK, rec.Code, "step %d: expected 200", idx)
		s.Equal(want, rec.Body.String(), "step %d: body", idx)
	}
}

func (s *TestSuite) TestAllDown() {
	//Arrange
	router := http.NewServeMux()
	backendPool := s.newBackendPool(s.urls)
	strategyRR := strategy.NewRoundRobinBalancer(backendPool)
	b := balancer.NewBalancer(strategyRR, backendPool, s.log)
	b.RegisterBalancer(router)

	backendPool.Backends[0].SetAlive(false)
	backendPool.Backends[1].SetAlive(false)
	backendPool.Backends[2].SetAlive(false)

	// Create Request
	ctx := context.WithValue(context.Background(), middleware.CtxKeyRequestID, "req-3")
	req := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	// Make Request
	router.ServeHTTP(rec, req)

	//Assert
	s.Equal(http.StatusServiceUnavailable, rec.Code)
	expectedBody := &model.ErrorResponse{
		Code:    503,
		Message: balancer.ServiceNotAvailable.Error(),
	}
	var body *model.ErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&body)
	s.Require().NoError(err, "decode error response")
	s.Equal(expectedBody.Message, body.Message)
}
