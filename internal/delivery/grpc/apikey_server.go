package grpcdelivery

import (
	"context"

	"github.com/user/auth/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/user/auth/pkg/authv1"
)

// APIKeyServer implements the APIKeyService gRPC interface.
type APIKeyServer struct {
	pb.UnimplementedAPIKeyServiceServer
	apikeyUsecase *usecase.APIKeyUsecase
}

// NewAPIKeyServer creates a new APIKeyServer.
func NewAPIKeyServer(apikeyUsecase *usecase.APIKeyUsecase) *APIKeyServer {
	return &APIKeyServer{apikeyUsecase: apikeyUsecase}
}

// CreateAPIKey creates a new API key.
func (s *APIKeyServer) CreateAPIKey(ctx context.Context, req *pb.CreateAPIKeyRequest) (*pb.CreateAPIKeyResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// ListAPIKeys returns a paginated list of API keys.
func (s *APIKeyServer) ListAPIKeys(ctx context.Context, req *pb.ListAPIKeysRequest) (*pb.ListAPIKeysResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// ValidateAPIKey validates an API key.
func (s *APIKeyServer) ValidateAPIKey(ctx context.Context, req *pb.ValidateAPIKeyRequest) (*pb.ValidateAPIKeyResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// RevokeAPIKey deactivates an API key.
func (s *APIKeyServer) RevokeAPIKey(ctx context.Context, req *pb.RevokeAPIKeyRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
