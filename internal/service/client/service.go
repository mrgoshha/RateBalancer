package client

import (
	"RateBalancer/internal/adapter/dbs"
	"RateBalancer/internal/model"
	"RateBalancer/internal/service"
	"RateBalancer/internal/service/limiter"
	"RateBalancer/pkg/hash"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Service struct {
	repository dbs.Client
	hasher     hash.Hasher
	cfgLimiter *limiter.Config
}

func NewService(r dbs.Client, h hash.Hasher, cl *limiter.Config) *Service {
	return &Service{
		repository: r,
		hasher:     h,
		cfgLimiter: cl,
	}
}

func (s *Service) Create(ctx context.Context, req *service.CreateClientRequest) (*service.ClientCredentials, error) {
	id := uuid.New().String()
	apiKey := uuid.New().String()
	apiKeyHash, err := s.hasher.Hash(apiKey)
	if err != nil {
		return nil, fmt.Errorf(`hash apikey %w`, err)
	}

	client, err := model.NewClient(id, apiKeyHash, req.Capacity, req.PerSecond, time.Now().UTC())
	if err != nil {
		return nil, fmt.Errorf(`create client %w`, err)
	}

	if err = s.repository.Create(ctx, client); err != nil {
		return nil, err
	}

	res := &service.ClientCredentials{
		Id:     id,
		ApiKey: apiKey,
	}

	return res, nil
}

func (s *Service) Get(ctx context.Context, id string) (*model.Client, error) {
	return s.repository.Get(ctx, id)
}

func (s *Service) Update(ctx context.Context, id string, req *service.UpdateClientRequest) (*model.Client, error) {
	client, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = client.UpdateLimits(req.Capacity, req.PerSecond)
	if err != nil {
		return nil, err
	}

	return s.repository.Update(ctx, client)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}
