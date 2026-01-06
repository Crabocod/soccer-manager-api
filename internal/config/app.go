package config

type AppConfig struct {
	Name        string `envconfig:"APP_NAME" default:"soccer-manager"`
	Environment string `envconfig:"APP_ENV" default:"development"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"info"`
}
