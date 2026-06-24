package pg

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	UserRepo     *UserRepo
	SessionRepo  *SessionRepo
	PasswordRepo *PasswordRepo
	RoleRepo     *RoleRepo
	APIKeyRepo   *APIKeyRepo
}

func NewRepositories(pool *pgxpool.Pool) *Repositories {
	dbPool := NewDBPool(pool)
	return &Repositories{
		UserRepo:     NewUserRepo(dbPool),
		SessionRepo:  NewSessionRepo(dbPool),
		PasswordRepo: NewPasswordRepo(dbPool),
		RoleRepo:     NewRoleRepo(dbPool),
		APIKeyRepo:   NewAPIKeyRepo(dbPool),
	}
}
