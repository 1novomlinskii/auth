// Package app contains the application lifecycle: initialization, run, and shutdown.
package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/user/auth/internal/config"
	"github.com/user/auth/internal/delivery/http"
)

// App represents the auth service application.
type App struct {
	cfg     *config.Config
	logger  zerolog.Logger
	pool    *pgxpool.Pool
	httpSrv *http.Server
}

// New creates a new App with the given database pool.
func New(cfg *config.Config, logger zerolog.Logger, pool *pgxpool.Pool) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
		pool:   pool,
	}
}

// Run starts the HTTP server and blocks until the context is cancelled,
// then performs a graceful shutdown.
func (a *App) Run(ctx context.Context) error {
	if err := a.startHTTPServer(ctx); err != nil {
		return fmt.Errorf("start http server: %w", err)
	}

	<-ctx.Done()
	a.logger.Info().Msg("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.httpSrv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}

	a.logger.Info().Msg("http server stopped")

	return nil
}

// startHTTPServer configures routes and starts the HTTP server in a goroutine.
func (a *App) startHTTPServer(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", httpdelivery.HealthHandler(a.pool, a.logger))

	handler := hlog.NewHandler(a.logger)(mux)
	handler = hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("request")
	})(handler)
	handler = hlog.RequestIDHandler("req_id", "X-Request-ID")(handler)

	addr := ":" + a.cfg.HTTP.Port
	a.httpSrv = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		a.logger.Info().Str("addr", addr).Msg("starting http server")
		if err := a.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal().Err(err).Msg("http server error")
		}
	}()

	return nil
}