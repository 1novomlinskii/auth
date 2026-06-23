// Package app contains the application lifecycle: initialization, run, and shutdown.
package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/user/auth/internal/config"
	grpcdelivery "github.com/user/auth/internal/delivery/grpc"
	httpdelivery "github.com/user/auth/internal/delivery/http"
	"github.com/user/auth/internal/repository/pg"
	"github.com/user/auth/internal/usecase"
	pb "github.com/user/auth/pkg/authv1"
)

// App represents the auth service application.
type App struct {
	cfg     *config.Config
	logger  zerolog.Logger
	pool    *pgxpool.Pool
	httpSrv *http.Server
	grpcSrv *grpc.Server
}

// New creates a new App with the given database pool.
func New(cfg *config.Config, logger zerolog.Logger, pool *pgxpool.Pool) *App {
	// Repositories
	dbPool := pg.NewDBPool(pool)
	userRepo := pg.NewUserRepo(dbPool)
	sessionRepo := pg.NewSessionRepo(dbPool)
	passwordRepo := pg.NewPasswordRepo(dbPool)
	roleRepo := pg.NewRoleRepo(dbPool)
	apikeyRepo := pg.NewAPIKeyRepo(dbPool)

	// Usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo, passwordRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)
	roleUsecase := usecase.NewRoleUsecase(roleRepo)
	apikeyUsecase := usecase.NewAPIKeyUsecase(apikeyRepo)

	// gRPC server
	grpcSrv := grpcdelivery.NewServer()
	pb.RegisterAuthServiceServer(grpcSrv, grpcdelivery.NewAuthServer(authUsecase))
	pb.RegisterUserServiceServer(grpcSrv, grpcdelivery.NewUserServer(userUsecase))
	pb.RegisterTokenServiceServer(grpcSrv, grpcdelivery.NewTokenServer(authUsecase))
	pb.RegisterRoleServiceServer(grpcSrv, grpcdelivery.NewRoleServer(roleUsecase))
	pb.RegisterAPIKeyServiceServer(grpcSrv, grpcdelivery.NewAPIKeyServer(apikeyUsecase))

	// Enable reflection for debugging (disable in production).
	reflection.Register(grpcSrv)

	// HTTP handlers
	authHandler := httpdelivery.NewAuthHandler(authUsecase)
	userHandler := httpdelivery.NewUserHandler(userUsecase)
	roleHandler := httpdelivery.NewRoleHandler(roleUsecase)
	apikeyHandler := httpdelivery.NewAPIKeyHandler(apikeyUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", httpdelivery.HealthHandler(pool, logger))

	// Auth routes
	mux.HandleFunc("POST /v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /v1/auth/refresh", authHandler.RefreshToken)
	mux.HandleFunc("POST /v1/auth/logout", authHandler.Logout)
	mux.HandleFunc("GET /v1/auth/oauth2/{provider}/url", authHandler.GetOAuth2URL)
	mux.HandleFunc("POST /v1/auth/oauth2/login", authHandler.LoginOAuth2)

	// User routes
	mux.HandleFunc("GET /v1/users/{user_id}", userHandler.GetUser)
	mux.HandleFunc("PATCH /v1/users/{user_id}", userHandler.UpdateUser)
	mux.HandleFunc("GET /v1/users", userHandler.ListUsers)

	// Role routes
	mux.HandleFunc("POST /v1/roles", roleHandler.CreateRole)
	mux.HandleFunc("GET /v1/roles/{role_id}", roleHandler.GetRole)
	mux.HandleFunc("GET /v1/roles", roleHandler.ListRoles)
	mux.HandleFunc("PATCH /v1/roles/{role_id}", roleHandler.UpdateRole)
	mux.HandleFunc("DELETE /v1/roles/{role_id}", roleHandler.DeleteRole)
	mux.HandleFunc("POST /v1/roles/{role_id}/assign", roleHandler.AssignRole)
	mux.HandleFunc("POST /v1/roles/{role_id}/revoke", roleHandler.RevokeRole)
	mux.HandleFunc("POST /v1/permissions/check", roleHandler.CheckPermission)
	mux.HandleFunc("GET /v1/users/{user_id}/roles", roleHandler.GetUserRoles)

	// API key routes
	mux.HandleFunc("POST /v1/api-keys", apikeyHandler.CreateAPIKey)
	mux.HandleFunc("GET /v1/api-keys", apikeyHandler.ListAPIKeys)
	mux.HandleFunc("POST /v1/api-keys/validate", apikeyHandler.ValidateAPIKey)
	mux.HandleFunc("POST /v1/api-keys/{api_key_id}/revoke", apikeyHandler.RevokeAPIKey)

	// Middleware chain.
	handler := httpdelivery.CORS(mux)
	handler = hlog.NewHandler(logger)(handler)
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

	return &App{
		cfg:    cfg,
		logger: logger,
		pool:   pool,
		httpSrv: &http.Server{
			Addr:         ":" + cfg.HTTP.Port,
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		grpcSrv: grpcSrv,
	}
}

// Run starts the gRPC and HTTP servers and blocks until the context is cancelled,
// then performs a graceful shutdown.
func (a *App) Run(ctx context.Context) error {
	// Start gRPC server.
	grpcAddr := ":" + a.cfg.GRPC.Port
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("listen grpc: %w", err)
	}

	go func() {
		a.logger.Info().Str("addr", grpcAddr).Msg("starting gRPC server")
		if err := a.grpcSrv.Serve(lis); err != nil {
			a.logger.Fatal().Err(err).Msg("gRPC server error")
		}
	}()

	// Start HTTP server.
	go func() {
		httpAddr := ":" + a.cfg.HTTP.Port
		a.logger.Info().Str("addr", httpAddr).Msg("starting HTTP server")
		if err := a.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal().Err(err).Msg("HTTP server error")
		}
	}()

	<-ctx.Done()
	a.logger.Info().Msg("shutting down...")

	// Graceful shutdown: gRPC first, then HTTP.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan struct{}, 1)

	go func() {
		a.grpcSrv.GracefulStop()
		a.logger.Info().Msg("gRPC server stopped")

		if err := a.httpSrv.Shutdown(shutdownCtx); err != nil {
			a.logger.Error().Err(err).Msg("HTTP server shutdown error")
		}
		a.logger.Info().Msg("HTTP server stopped")

		close(done)
	}()

	select {
	case <-done:
	case <-shutdownCtx.Done():
		a.logger.Warn().Msg("shutdown timed out, forcing stop")
		a.grpcSrv.Stop()
	}

	return nil
}
