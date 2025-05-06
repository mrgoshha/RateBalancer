package balancer

import (
	"fmt"
	"net/url"
	"time"
)

type BackendConfig struct {
	URL    *url.URL `yaml:"-"`
	RawURL string   `yaml:"url"`
}

type Config struct {
	Strategy           string          `yaml:"strategy"`
	Backends           []BackendConfig `yaml:"backends"`
	HealthyThreshold   int             `yaml:"healthy_threshold"`
	UnhealthyThreshold int             `yaml:"unhealthy_threshold"`
	PingTimeout        time.Duration   `yaml:"timeout"`
}

func (lbc *Config) ParseBackends() error {
	if len(lbc.Backends) == 0 {
		return fmt.Errorf("no backends configured")
	}

	for i := range lbc.Backends {
		raw := lbc.Backends[i].RawURL
		u, err := url.Parse(raw)
		if err != nil {
			return fmt.Errorf("invalid backend URL %q: %w", raw, err)
		}
		lbc.Backends[i].URL = u
	}
	return nil
}
