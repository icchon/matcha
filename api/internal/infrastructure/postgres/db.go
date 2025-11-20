package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var _ DBTX = (*sqlx.DB)(nil)
var _ DBTX = (*sqlx.Tx)(nil)

type DBTX interface {
	sqlx.ExtContext

	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}
