package limiter

type Config struct {
	Capacity  int64 `yaml:"default_capacity"`
	PerSecond int64 `yaml:"default_rate_per_sec"`
}
