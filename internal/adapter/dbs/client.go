package dbs

import (
	"RateBalancer/internal/model"
	"context"
)

type Client interface {
	Create(ctx context.Context, client *model.Client) error
	Get(ctx context.Context, id string) (*model.Client, error)
	Update(ctx context.Context, client *model.Client) (*model.Client, error)
	Delete(ctx context.Context, id string) error
}
