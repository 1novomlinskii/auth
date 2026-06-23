package pg

import (
	"context"
	"errors"

	"github.com/user/auth/internal/entity"
	"github.com/user/auth/internal/usecase"
)

// Ensure *SessionRepo implements usecase.SessionRepository.
var _ usecase.SessionRepository = (*SessionRepo)(nil)

// SessionRepo implements the usecase.SessionRepository interface.
type SessionRepo struct {
	pool *DBPool
}

// NewSessionRepo creates a new SessionRepo.
func NewSessionRepo(pool *DBPool) *SessionRepo {
	return &SessionRepo{pool: pool}
}

// Create inserts a new session record.
func (r *SessionRepo) Create(ctx context.Context, session *entity.Session) error {
	return errors.New("not implemented")
}

// GetByRefreshToken retrieves a session by its refresh token hash.
func (r *SessionRepo) GetByRefreshToken(ctx context.Context, refreshTokenHash string) (*entity.Session, error) {
	return nil, errors.New("not implemented")
}

// Revoke marks a session as revoked.
func (r *SessionRepo) Revoke(ctx context.Context, sessionID string) error {
	return errors.New("not implemented")
}

// RevokeAllUserSessions revokes all active sessions for a user.
func (r *SessionRepo) RevokeAllUserSessions(ctx context.Context, userID string) error {
	return errors.New("not implemented")
}

// Ensure *PasswordRepo implements usecase.PasswordRepository.
var _ usecase.PasswordRepository = (*PasswordRepo)(nil)

// PasswordRepo implements the usecase.PasswordRepository interface.
type PasswordRepo struct {
	pool *DBPool
}

// NewPasswordRepo creates a new PasswordRepo.
func NewPasswordRepo(pool *DBPool) *PasswordRepo {
	return &PasswordRepo{pool: pool}
}

// Set stores a password hash for a user.
func (r *PasswordRepo) Set(ctx context.Context, userID, passwordHash string) error {
	return errors.New("not implemented")
}

// GetByUserID retrieves the password hash for a user.
func (r *PasswordRepo) GetByUserID(ctx context.Context, userID string) (string, error) {
	return "", errors.New("not implemented")
}

// Update changes the password hash for a user.
func (r *PasswordRepo) Update(ctx context.Context, userID, passwordHash string) error {
	return errors.New("not implemented")
}
