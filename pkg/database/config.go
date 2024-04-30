package database

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Environment      string `envconfig:"env"`
	DatabaseHost     string `envconfig:"database_host"`
	DatabasePort     string `envconfig:"database_port"`
	DatabaseDBName   string `envconfig:"database_name"`
	DatabaseSchema   string `envconfig:"database_schema"`
	DatabaseUsername string `envconfig:"database_username"`
	DatabasePassword string `envconfig:"database_password"`
	MigrationsDir    string `envconfig:"migrations_dir"`
}

func NewConfig() (cfg Config, err error) {
	err = envconfig.Process("", &cfg)
	return
}
