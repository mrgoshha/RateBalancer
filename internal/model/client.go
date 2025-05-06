package model

import (
	"errors"
	"time"
)

var (
	InvalidLimits = errors.New("limits is not valid")
)

type Client struct {
	Id         string
	ApiKey     string
	Tokens     int
	LastRefill time.Time
	Capacity   *int
	PerSecond  *int
}

func NewClient(id, apiKey string, capacity, perSecond *int, now time.Time) (*Client, error) {
	c := &Client{
		Id:         id,
		ApiKey:     apiKey,
		LastRefill: now,
		Capacity:   capacity,
		PerSecond:  perSecond,
	}

	if capacity != nil {
		c.Tokens = *capacity
	}

	if err := c.validateLimits(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) UpdateLimits(capacity, perSecond *int) error {
	if capacity != nil {
		c.Capacity = capacity
	}
	if perSecond != nil {
		c.PerSecond = perSecond
	}
	return c.validateLimits()
}

func (c *Client) validateLimits() error {
	switch {
	case c.Capacity == nil && c.PerSecond == nil:
		return nil
	case c.Capacity != nil && c.PerSecond != nil && *c.Capacity >= 0 && *c.PerSecond > 0:
		return nil
	default:
		return InvalidLimits
	}
}
