package adminserver

import (
	apiModel "RateBalancer/internal/handler/http/model"
	"RateBalancer/internal/model"
	"RateBalancer/internal/service"
)

func ToCreateClientRequest(c *apiModel.Client) *service.CreateClientRequest {
	return &service.CreateClientRequest{
		Id:        c.Id,
		Capacity:  c.Capacity,
		PerSecond: c.PerSecond,
	}
}

func ToClientCredentialsApiModel(c *service.ClientCredentials) *apiModel.ClientCredentials {
	return &apiModel.ClientCredentials{
		Id:     c.Id,
		ApiKey: c.ApiKey,
	}
}

func ToUpdateClientRequest(c *apiModel.UpdateClient) *service.UpdateClientRequest {
	return &service.UpdateClientRequest{
		Capacity:  c.Capacity,
		PerSecond: c.PerSecond,
	}
}

func ToClientApiModel(c *model.Client) *apiModel.Client {
	return &apiModel.Client{
		Id:        c.Id,
		Capacity:  c.Capacity,
		PerSecond: c.PerSecond,
	}
}
