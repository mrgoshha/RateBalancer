package balancer

import "errors"

var (
	ServiceNotAvailable = errors.New("service not available")
	ServerNotExist      = errors.New("server not exist")
)
