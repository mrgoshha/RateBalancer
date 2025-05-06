package model

import "time"

type GetClient struct {
	Id         string    `json:"id"`
	Tokens     int       `json:"tokens"`
	LastRefill time.Time `json:"last_refill"`
	Capacity   int       `json:"capacity,omitempty"`
	PerSecond  int       `json:"per_second,omitempty"`
}
