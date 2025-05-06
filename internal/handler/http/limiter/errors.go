package limiter

import "errors"

var (
	NoApiKey          = errors.New("no api key")
	RateLimitExceeded = errors.New("rate limit exceeded")
)
