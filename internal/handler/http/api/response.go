package api

import (
	"RateBalancer/internal/adapter/dbs"
	"RateBalancer/internal/handler/http/middleware"
	apimodel "RateBalancer/internal/handler/http/model"
	"RateBalancer/internal/model"
	"RateBalancer/internal/service/limiter"
	"encoding/json"
	"errors"
	"net/http"
)

func ErrorResponseWithCode(w http.ResponseWriter, r *http.Request, code int, err error) {
	errorResponse := apimodel.ErrorResponse{
		Code:      code,
		RequestId: middleware.GetRequestIdFromContext(r),
	}
	if code != http.StatusInternalServerError {
		errorResponse.Message = err.Error()
	}
	response(w, r, code, errorResponse)

}

func ErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	var code int

	switch {
	case errors.Is(err, dbs.ErrorRecordNotFound):
		code = http.StatusNotFound
	case errors.Is(err, dbs.ErrorRecordAlreadyExists):
		code = http.StatusConflict
	case errors.Is(err, model.InvalidLimits):
		code = http.StatusBadRequest
	case errors.Is(err, limiter.InvalidAPIKey):
		code = http.StatusUnauthorized
	default:
		code = http.StatusInternalServerError
	}
	ErrorResponseWithCode(w, r, code, err)
}

func response(w http.ResponseWriter, _ *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
