// Package pg implements repository interfaces using PostgreSQL.
package pg

//go:generate go install github.com/valyala/quicktemplate/qtc@v1.7.0
//go:generate qtc -dir=query

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBPool is a shared database pool used by all repository types.
type DBPool struct {
	pool *pgxpool.Pool
}

// NewDBPool creates a new DBPool wrapper.
func NewDBPool(pool *pgxpool.Pool) *DBPool {
	return &DBPool{pool: pool}
}
