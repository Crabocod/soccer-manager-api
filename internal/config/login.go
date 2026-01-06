package config

import "time"

type LoginConfig struct {
	MaxLoginAttempts      int           `envconfig:"LOGIN_MAX_ATTEMPTS" default:"5"`
	LoginAttemptTTL       time.Duration `envconfig:"LOGIN_ATTEMPT_TTL" default:"15m"`
	TeamCacheTTL          time.Duration `envconfig:"TEAM_CACHE_TTL" default:"5m"`
}
