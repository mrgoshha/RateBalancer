package api

import (
	"RateBalancer/internal/handler/http/adminserver"
	"RateBalancer/internal/handler/http/middleware"
	"RateBalancer/internal/handler/http/model"
	"RateBalancer/internal/service"
	"encoding/json"
	"log/slog"
	"net/http"
)

type ClientController struct {
	log     *slog.Logger
	service service.Client
}

func NewClientController(log *slog.Logger, service service.Client, router *http.ServeMux) *ClientController {
	c := &ClientController{
		log:     log,
		service: service,
	}

	router.HandleFunc("POST /client", c.Create)
	router.HandleFunc("GET /client/{id}", c.Get)
	router.HandleFunc("PATCH /client/{id}", c.Update)
	router.HandleFunc("DELETE /client/{id}", c.Delete)

	return c
}

func (c *ClientController) Create(w http.ResponseWriter, r *http.Request) {
	const op = "handler.http.api.Client"

	log := c.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetRequestIdFromContext(r)),
	)
	req := &model.Client{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Error("failed to decode request")
		ErrorResponseWithCode(w, r, http.StatusBadRequest, err)
		return
	}

	credentials, err := c.service.Create(r.Context(), adminserver.ToCreateClientRequest(req))
	if err != nil {
		log.Error("failed to create user", slog.String("error", err.Error()))
		ErrorResponse(w, r, err)
		return
	}

	log.Info("create new client")

	res := adminserver.ToClientCredentialsApiModel(credentials)

	response(w, r, http.StatusOK, res)
}

func (c *ClientController) Get(w http.ResponseWriter, r *http.Request) {
	const op = "handler.http.api.UpdateClient"

	log := c.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetRequestIdFromContext(r)),
	)

	id := r.PathValue("id")
	client, err := c.service.Get(r.Context(), id)
	if err != nil {
		log.Error("failed to get user", slog.String("error", err.Error()))
		ErrorResponse(w, r, err)
		return
	}

	res := adminserver.ToClientApiModel(client)

	response(w, r, http.StatusOK, res)
}

func (c *ClientController) Update(w http.ResponseWriter, r *http.Request) {
	const op = "handler.http.api.UpdateClient"

	log := c.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetRequestIdFromContext(r)),
	)
	req := &model.UpdateClient{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Error("failed to decode request")
		ErrorResponseWithCode(w, r, http.StatusBadRequest, err)
		return
	}

	id := r.PathValue("id")

	client, err := c.service.Update(r.Context(), id, adminserver.ToUpdateClientRequest(req))
	if err != nil {
		log.Error("failed to update user", slog.String("error", err.Error()))
		ErrorResponse(w, r, err)
		return
	}

	log.Info("update client")

	res := adminserver.ToClientApiModel(client)

	response(w, r, http.StatusOK, res)
}

func (c *ClientController) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "handler.http.api.UpdateClient"

	log := c.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetRequestIdFromContext(r)),
	)

	id := r.PathValue("id")

	err := c.service.Delete(r.Context(), id)
	if err != nil {
		log.Error("failed to delete user", slog.String("error", err.Error()))
		ErrorResponse(w, r, err)
		return
	}

	response(w, r, http.StatusOK, nil)
}
