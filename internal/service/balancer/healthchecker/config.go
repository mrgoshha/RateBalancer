package healthchecker

import "time"

type Config struct {
	PingInterval time.Duration `yaml:"ping_interval"`
}
