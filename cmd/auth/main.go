package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"

	"github.com/user/auth/internal/app"
	"github.com/user/auth/internal/config"
)

func main() {
	// 1. Load configuration.
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "./config.yml"
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialise zerolog.
	logLevel, err := zerolog.ParseLevel(cfg.Log.Level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	var logger zerolog.Logger
	if cfg.Log.Pretty {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
			With().Timestamp().Caller().Logger()
	} else {
		logger = zerolog.New(os.Stdout).
			With().Timestamp().Caller().Logger()
	}

	// 3. Create context that cancels on SIGTERM / SIGINT.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// 4. Create database connection pool.
	poolCfg, err := pgxpool.ParseConfig(cfg.DB.URL)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse pool config")
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create pool")
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatal().Err(err).Msg("failed to ping database")
	}

	logger.Info().Msg("connected to database")

	// 5. Run database migrations.
	db := stdlib.OpenDBFromPool(pool)
	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		db,
		os.DirFS("migrations"),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create goose provider")
	}

	if _, err := provider.Up(ctx); err != nil {
		logger.Fatal().Err(err).Msg("failed to run migrations")
	}
	if err := db.Close(); err != nil {
		logger.Error().Err(err).Msg("failed to close migration db")
	}

	logger.Info().Msg("migrations completed")

	// 6. Create and run the application.
	application := app.New(cfg, logger, pool)

	logger.Info().Str("port", cfg.HTTP.Port).Msg("starting auth service")

	if err := application.Run(ctx); err != nil {
		logger.Fatal().Err(err).Msg("app exited with error")
	}

	logger.Info().Msg("auth service stopped gracefully")
}
