package grpcdelivery

import (
	"context"

	"github.com/user/auth/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/user/auth/pkg/authv1"
)

// RoleServer implements the RoleService gRPC interface.
type RoleServer struct {
	pb.UnimplementedRoleServiceServer
	roleUsecase *usecase.RoleUsecase
}

// NewRoleServer creates a new RoleServer.
func NewRoleServer(roleUsecase *usecase.RoleUsecase) *RoleServer {
	return &RoleServer{roleUsecase: roleUsecase}
}

// CreateRole creates a new role.
func (s *RoleServer) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.Role, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// GetRole returns a role by ID.
func (s *RoleServer) GetRole(ctx context.Context, req *pb.GetRoleRequest) (*pb.Role, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// ListRoles returns a paginated list of roles.
func (s *RoleServer) ListRoles(ctx context.Context, req *pb.ListRolesRequest) (*pb.ListRolesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// UpdateRole updates a role.
func (s *RoleServer) UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.Role, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// DeleteRole deletes a role.
func (s *RoleServer) DeleteRole(ctx context.Context, req *pb.DeleteRoleRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// AssignRole assigns a role to a user.
func (s *RoleServer) AssignRole(ctx context.Context, req *pb.AssignRoleRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// RevokeRole revokes a role from a user.
func (s *RoleServer) RevokeRole(ctx context.Context, req *pb.RevokeRoleRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// CheckPermission checks if a user has a specific permission.
func (s *RoleServer) CheckPermission(ctx context.Context, req *pb.CheckPermissionRequest) (*pb.CheckPermissionResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// GetUserRoles returns all roles for a user.
func (s *RoleServer) GetUserRoles(ctx context.Context, req *pb.GetUserRolesRequest) (*pb.GetUserRolesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
