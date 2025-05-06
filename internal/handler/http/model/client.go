package model

type Client struct {
	Id        string `json:"id"`
	Capacity  *int   `json:"capacity,omitempty"`
	PerSecond *int   `json:"per_second,omitempty"`
}
