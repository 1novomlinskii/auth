package usecase

import (
	"context"
	"errors"

	"github.com/user/auth/internal/entity"
)

// UserUsecase implements user management business logic.
type UserUsecase struct {
	repo UserRepository
}

// NewUserUsecase creates a new UserUsecase.
func NewUserUsecase(repo UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// Get returns a user by ID.
func (uc *UserUsecase) Get(ctx context.Context, userID string) (*entity.User, error) {
	return nil, errors.New("not implemented")
}

// Update updates a user's profile fields.
func (uc *UserUsecase) Update(ctx context.Context, userID, displayName, locale, timezone string) (*entity.User, error) {
	return nil, errors.New("not implemented")
}

// List returns a paginated list of users.
func (uc *UserUsecase) List(ctx context.Context, page, pageSize int) ([]*entity.User, int, error) {
	return nil, 0, errors.New("not implemented")
}
