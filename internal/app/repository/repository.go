package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type connection interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Repository struct {
	conn connection
}

func NewRepository(conn connection) *Repository {
	return &Repository{
		conn: conn,
	}
}
