package pg

import (
	"context"
	"errors"

	"github.com/user/auth/internal/entity"
	"github.com/user/auth/internal/usecase"
)

// Ensure *UserRepo implements usecase.UserRepository.
var _ usecase.UserRepository = (*UserRepo)(nil)

// UserRepo implements the usecase.UserRepository interface.
type UserRepo struct {
	pool *DBPool
}

// NewUserRepo creates a new UserRepo.
func NewUserRepo(pool *DBPool) *UserRepo {
	return &UserRepo{pool: pool}
}

// Create inserts a new user into the database.
func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	return errors.New("not implemented")
}

// GetByID retrieves a user by their primary key.
func (r *UserRepo) GetByID(ctx context.Context, userID string) (*entity.User, error) {
	return nil, errors.New("not implemented")
}

// GetByLogin retrieves a user by their login name.
func (r *UserRepo) GetByLogin(ctx context.Context, login string) (*entity.User, error) {
	return nil, errors.New("not implemented")
}

// GetByEmail retrieves a user by their primary email.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, errors.New("not implemented")
}

// Update modifies an existing user record.
func (r *UserRepo) Update(ctx context.Context, user *entity.User) error {
	return errors.New("not implemented")
}

// List returns a paginated list of users and the total count.
func (r *UserRepo) List(ctx context.Context, page, pageSize int) ([]*entity.User, int, error) {
	return nil, 0, errors.New("not implemented")
}
