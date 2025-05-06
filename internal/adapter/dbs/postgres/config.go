package postgres

type Config struct {
	DbName   string `yaml:"postgres_db"`
	Host     string `yaml:"postgres_host"`
	Port     string `yaml:"postgres_ports"`
	User     string `yaml:"postgres_user"`
	Password string `yaml:"postgres_password"`
}
