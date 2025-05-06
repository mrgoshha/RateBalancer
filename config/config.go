package config

import (
	"RateBalancer/internal/adapter/dbs/postgres"
	"RateBalancer/internal/handler/http/adminserver"
	"RateBalancer/internal/handler/http/server"
	balancer2 "RateBalancer/internal/service/balancer"
	"RateBalancer/internal/service/balancer/healthchecker"
	"RateBalancer/internal/service/limiter"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server        *server.Config        `yaml:"server"`
	AdminServer   *adminserver.Config   `yaml:"adminServer"`
	Database      *postgres.Config      `yaml:"database"`
	Limiter       *limiter.Config       `yaml:"rateLimiter"`
	LoadBalancer  *balancer2.Config     `yaml:"loadBalancer"`
	HealthChecker *healthchecker.Config `yaml:"healthChecker"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := cfg.LoadBalancer.ParseBackends(); err != nil {
		return nil, fmt.Errorf("loadBalancer parse backends: %w", err)
	}

	return &cfg, nil
}
