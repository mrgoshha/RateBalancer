package service

import "RateBalancer/internal/service/balancer"

const (
	RoundRobin string = "round_robin"
	Random     string = "random"
)

type Strategy interface {
	GetNext() (*balancer.Backend, error)
}
