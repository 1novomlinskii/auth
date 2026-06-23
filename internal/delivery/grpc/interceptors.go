package grpcdelivery

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RecoveryInterceptor recovers from panics in gRPC handlers.
func RecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		defer func() {
			if rec := recover(); rec != nil {
				// TODO: log the panic
			}
		}()
		return handler(ctx, req)
	}
}

// LoggingInterceptor logs gRPC requests.
func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// TODO: log method name, duration, etc.
		return handler(ctx, req)
	}
}

// AuthInterceptor checks JWT tokens on protected gRPC methods.
// Currently a stub that passes through all requests.
func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// TODO: extract token from metadata, validate, inject claims into context
		return handler(ctx, req)
	}
}

// Unimplemented is a helper error for stubs.
var errUnimplemented = status.Error(codes.Unimplemented, "not implemented")

// NewServer creates a gRPC server with interceptors configured.
func NewServer() *grpc.Server {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			RecoveryInterceptor(),
			LoggingInterceptor(),
			AuthInterceptor(),
		),
	)
	return srv
}
