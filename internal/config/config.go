package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Server   ServerConfig
	Login    LoginConfig
}

func GetConfig() (*Config, error) {
	var conf Config

	if err := envconfig.Process("", &conf); err != nil {
		return nil, fmt.Errorf("read config from env vars: %w", err)
	}

	return &conf, nil
}
