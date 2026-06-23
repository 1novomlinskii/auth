package grpcdelivery

import (
	"context"

	"github.com/user/auth/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/user/auth/pkg/authv1"
)

// AuthServer implements the AuthService gRPC interface.
type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	authUsecase *usecase.AuthUsecase
}

// NewAuthServer creates a new AuthServer.
func NewAuthServer(authUsecase *usecase.AuthUsecase) *AuthServer {
	return &AuthServer{authUsecase: authUsecase}
}

// Register creates a new user account.
func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// Login authenticates a user by login/email and password.
func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// GetOAuth2URL returns an OAuth2 authorization URL for the given provider.
func (s *AuthServer) GetOAuth2URL(ctx context.Context, req *pb.GetOAuth2URLRequest) (*pb.GetOAuth2URLResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// LoginOAuth2 authenticates a user using an OAuth2 authorization code.
func (s *AuthServer) LoginOAuth2(ctx context.Context, req *pb.LoginOAuth2Request) (*pb.AuthResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// RefreshToken issues a new access token.
func (s *AuthServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.AuthResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// Logout revokes the given refresh token.
func (s *AuthServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
