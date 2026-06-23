package pg

import (
	"context"
	"errors"

	"github.com/user/auth/internal/entity"
	"github.com/user/auth/internal/usecase"
)

// Ensure *RoleRepo implements usecase.RoleRepository.
var _ usecase.RoleRepository = (*RoleRepo)(nil)

// RoleRepo implements the usecase.RoleRepository interface.
type RoleRepo struct {
	pool *DBPool
}

// NewRoleRepo creates a new RoleRepo.
func NewRoleRepo(pool *DBPool) *RoleRepo {
	return &RoleRepo{pool: pool}
}

// Create inserts a new role.
func (r *RoleRepo) Create(ctx context.Context, role *entity.Role) error {
	return errors.New("not implemented")
}

// GetByID retrieves a role by ID.
func (r *RoleRepo) GetByID(ctx context.Context, roleID string) (*entity.Role, error) {
	return nil, errors.New("not implemented")
}

// List returns a paginated list of roles.
func (r *RoleRepo) List(ctx context.Context, page, pageSize int) ([]*entity.Role, int, error) {
	return nil, 0, errors.New("not implemented")
}

// Update modifies an existing role.
func (r *RoleRepo) Update(ctx context.Context, role *entity.Role) error {
	return errors.New("not implemented")
}

// Delete removes a role by ID.
func (r *RoleRepo) Delete(ctx context.Context, roleID string) error {
	return errors.New("not implemented")
}

// AssignToUser grants a role to a user.
func (r *RoleRepo) AssignToUser(ctx context.Context, userID, roleID string) error {
	return errors.New("not implemented")
}

// RevokeFromUser removes a role from a user.
func (r *RoleRepo) RevokeFromUser(ctx context.Context, userID, roleID string) error {
	return errors.New("not implemented")
}

// GetUserRoles returns all roles for a given user.
func (r *RoleRepo) GetUserRoles(ctx context.Context, userID string) ([]*entity.Role, error) {
	return nil, errors.New("not implemented")
}

// CheckPermission verifies if a user has a specific permission.
func (r *RoleRepo) CheckPermission(ctx context.Context, userID, permission string) (bool, error) {
	return false, errors.New("not implemented")
}
