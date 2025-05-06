package service

import "context"

type Limiter interface {
	ConsumeTokens(ctx context.Context, apiKey string) (bool, error)
}
