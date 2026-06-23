package pg

import (
	"context"
	"errors"

	"github.com/user/auth/internal/entity"
	"github.com/user/auth/internal/usecase"
)

// Ensure *APIKeyRepo implements usecase.APIKeyRepository.
var _ usecase.APIKeyRepository = (*APIKeyRepo)(nil)

// APIKeyRepo implements the usecase.APIKeyRepository interface.
type APIKeyRepo struct {
	pool *DBPool
}

// NewAPIKeyRepo creates a new APIKeyRepo.
func NewAPIKeyRepo(pool *DBPool) *APIKeyRepo {
	return &APIKeyRepo{pool: pool}
}

// Create persists a new API key record.
func (r *APIKeyRepo) Create(ctx context.Context, key *entity.APIKey) error {
	return errors.New("not implemented")
}

// GetByHash retrieves an API key by its hash.
func (r *APIKeyRepo) GetByHash(ctx context.Context, keyHash string) (*entity.APIKey, error) {
	return nil, errors.New("not implemented")
}

// List returns a paginated list of API keys.
func (r *APIKeyRepo) List(ctx context.Context, page, pageSize int) ([]*entity.APIKey, int, error) {
	return nil, 0, errors.New("not implemented")
}

// Revoke deactivates an API key by ID.
func (r *APIKeyRepo) Revoke(ctx context.Context, apiKeyID string) error {
	return errors.New("not implemented")
}
