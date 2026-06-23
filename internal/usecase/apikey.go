package usecase

import (
	"context"
	"errors"

	"github.com/user/auth/internal/entity"
)

// APIKeyRepository defines the interface for API key persistence.
type APIKeyRepository interface {
	Create(ctx context.Context, key *entity.APIKey) error
	GetByHash(ctx context.Context, keyHash string) (*entity.APIKey, error)
	List(ctx context.Context, page, pageSize int) ([]*entity.APIKey, int, error)
	Revoke(ctx context.Context, apiKeyID string) error
}

// APIKeyUsecase implements API key management business logic.
type APIKeyUsecase struct {
	repo APIKeyRepository
}

// NewAPIKeyUsecase creates a new APIKeyUsecase.
func NewAPIKeyUsecase(repo APIKeyRepository) *APIKeyUsecase {
	return &APIKeyUsecase{repo: repo}
}

// Create generates a new API key.
func (uc *APIKeyUsecase) Create(ctx context.Context, serviceName string, permissions []string, expiresAt *string) (*entity.APIKey, string, error) {
	return nil, "", errors.New("not implemented")
}

// Validate checks if an API key is valid and returns its metadata.
func (uc *APIKeyUsecase) Validate(ctx context.Context, apiKey string) (*entity.APIKey, error) {
	return nil, errors.New("not implemented")
}

// Revoke deactivates an API key.
func (uc *APIKeyUsecase) Revoke(ctx context.Context, apiKeyID string) error {
	return errors.New("not implemented")
}

// List returns a paginated list of API keys.
func (uc *APIKeyUsecase) List(ctx context.Context, page, pageSize int) ([]*entity.APIKey, int, error) {
	return nil, 0, errors.New("not implemented")
}
