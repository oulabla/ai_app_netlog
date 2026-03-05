package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type connection interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type Repository struct {
	conn connection
}

func NewRepository(conn connection) *Repository {
	return &Repository{
		conn: conn,
	}
}
