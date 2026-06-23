package grpcdelivery

import (
	"context"

	"github.com/user/auth/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/user/auth/pkg/authv1"
)

// UserServer implements the UserService gRPC interface.
type UserServer struct {
	pb.UnimplementedUserServiceServer
	userUsecase *usecase.UserUsecase
}

// NewUserServer creates a new UserServer.
func NewUserServer(userUsecase *usecase.UserUsecase) *UserServer {
	return &UserServer{userUsecase: userUsecase}
}

// GetUser returns a user by ID.
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// UpdateUser updates a user's profile.
func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// ListUsers returns a paginated list of users.
func (s *UserServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// TokenServer implements the TokenService gRPC interface.
type TokenServer struct {
	pb.UnimplementedTokenServiceServer
	authUsecase *usecase.AuthUsecase
}

// NewTokenServer creates a new TokenServer.
func NewTokenServer(authUsecase *usecase.AuthUsecase) *TokenServer {
	return &TokenServer{authUsecase: authUsecase}
}

// ValidateToken validates an access token and returns basic claims.
func (s *TokenServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// IntrospectToken returns detailed token information.
func (s *TokenServer) IntrospectToken(ctx context.Context, req *pb.IntrospectTokenRequest) (*pb.IntrospectTokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
