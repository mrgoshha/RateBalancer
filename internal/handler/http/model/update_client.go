package model

type UpdateClient struct {
	Capacity  *int `json:"capacity,omitempty"`
	PerSecond *int `json:"per_second,omitempty"`
}
