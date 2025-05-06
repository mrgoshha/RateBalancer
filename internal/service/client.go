package service

import (
	"RateBalancer/internal/model"
	"context"
)

type Client interface {
	Create(ctx context.Context, client *CreateClientRequest) (*ClientCredentials, error)
	Get(ctx context.Context, id string) (*model.Client, error)
	Update(ctx context.Context, id string, req *UpdateClientRequest) (*model.Client, error)
	Delete(ctx context.Context, id string) error
}

type CreateClientRequest struct {
	Id        string
	Capacity  *int
	PerSecond *int
}

type UpdateClientRequest struct {
	Capacity  *int
	PerSecond *int
}

type ClientCredentials struct {
	Id     string
	ApiKey string
}
