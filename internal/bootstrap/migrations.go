package bootstrap

import (
	"database/sql"
	"soccer_manager_service/migrations"

	"soccer_manager_service/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

func runMigrations(config *config.Config, logger *zap.Logger) error {
	logger.Info("running database migrations")

	db, err := sql.Open("pgx", config.Database.DSN())
	if err != nil {
		logger.Error("failed to open database for migrations", zap.Error(err))

		return err
	}
	defer db.Close()

	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error("failed to set goose dialect", zap.Error(err))

		return err
	}

	if err := goose.Up(db, "."); err != nil {
		logger.Error("failed to run migrations", zap.Error(err))

		return err
	}

	logger.Info("migrations completed successfully")

	return nil
}
