package bootstrap

import (
	"time"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

const (
	postgresBreakerName = "postgres"
	redisBreakerName    = "redis"
)

func initBreakers(logger *zap.Logger) (map[string]*gobreaker.CircuitBreaker, error) {
	mapBreakers := make(map[string]*gobreaker.CircuitBreaker)

	postgresBreaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        postgresBreakerName,
		MaxRequests: 3,
		Interval:    1 * time.Minute,
		Timeout:     15 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)

			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Info("circuit breaker state changed",
				zap.String("breaker", name),
				zap.String("from", from.String()),
				zap.String("to", to.String()))
		},
	})

	mapBreakers[postgresBreakerName] = postgresBreaker

	redisBreaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        redisBreakerName,
		MaxRequests: 3,
		Interval:    1 * time.Minute,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)

			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Info("circuit breaker state changed",
				zap.String("breaker", name),
				zap.String("from", from.String()),
				zap.String("to", to.String()))
		},
	})

	mapBreakers[redisBreakerName] = redisBreaker

	return mapBreakers, nil
}
