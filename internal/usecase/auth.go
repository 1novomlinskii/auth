// Package usecase contains application business logic.
package usecase

import (
	"context"
	"errors"

	"github.com/user/auth/internal/entity"
)

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, userID string) (*entity.User, error)
	GetByLogin(ctx context.Context, login string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
}

// SessionRepository defines the interface for session persistence.
type SessionRepository interface {
	Create(ctx context.Context, session *entity.Session) error
	GetByRefreshToken(ctx context.Context, refreshTokenHash string) (*entity.Session, error)
	Revoke(ctx context.Context, sessionID string) error
	RevokeAllUserSessions(ctx context.Context, userID string) error
}

// PasswordRepository defines the interface for password hash persistence.
type PasswordRepository interface {
	Set(ctx context.Context, userID, passwordHash string) error
	GetByUserID(ctx context.Context, userID string) (string, error)
	Update(ctx context.Context, userID, passwordHash string) error
}

// AuthResult contains tokens and user data returned by auth operations.
type AuthResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	TokenType    string
	User         *entity.User
}

// TokenClaims contains parsed JWT claims.
type TokenClaims struct {
	UserID   string
	TenantID string
	Roles    []string
}

// AuthUsecase implements authentication business logic.
type AuthUsecase struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
	passRepo    PasswordRepository
}

// NewAuthUsecase creates a new AuthUsecase.
func NewAuthUsecase(u UserRepository, s SessionRepository, p PasswordRepository) *AuthUsecase {
	return &AuthUsecase{userRepo: u, sessionRepo: s, passRepo: p}
}

// Register creates a new user account.
func (uc *AuthUsecase) Register(ctx context.Context, login, email, password string) (*AuthResult, error) {
	return nil, errors.New("not implemented")
}

// Login authenticates a user by login/email and password.
func (uc *AuthUsecase) Login(ctx context.Context, login, password string) (*AuthResult, error) {
	return nil, errors.New("not implemented")
}

// RefreshToken issues a new access token using a valid refresh token.
func (uc *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (*AuthResult, error) {
	return nil, errors.New("not implemented")
}

// Logout revokes the given refresh token.
func (uc *AuthUsecase) Logout(ctx context.Context, refreshToken string) error {
	return errors.New("not implemented")
}

// ValidateToken parses and validates an access token.
func (uc *AuthUsecase) ValidateToken(ctx context.Context, token string) (*TokenClaims, error) {
	return nil, errors.New("not implemented")
}
