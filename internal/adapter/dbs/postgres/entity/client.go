package entity

import "time"

type Client struct {
	Id         string    `db:"id"`
	ApiKey     string    `db:"api_key"`
	Tokens     int       `db:"tokens"`
	LastRefill time.Time `db:"last_refill"`
	Capacity   *int      `db:"capacity"`
	PerSecond  *int      `db:"per_second"`
}
