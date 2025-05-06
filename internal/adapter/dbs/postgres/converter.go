package postgres

import (
	"RateBalancer/internal/adapter/dbs/postgres/entity"
	"RateBalancer/internal/model"
)

func ToClientServiceModel(u *entity.Client) *model.Client {
	return &model.Client{
		Id:         u.Id,
		ApiKey:     u.ApiKey,
		Tokens:     u.Tokens,
		LastRefill: u.LastRefill,
		Capacity:   u.Capacity,
		PerSecond:  u.PerSecond,
	}
}
