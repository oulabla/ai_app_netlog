package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/oulabla/ai_app_netlog/internal/datastruct"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (*datastruct.Netlog, error) {
	query := `
		SELECT
			id,
			created_at,
			keywords,
			parameters,
			error,
			num_before_ai_filter,
			num_after_ai_filter,
			result_before_ai_filter,
			result,
			client_id,
			app_name
		FROM netlog
		WHERE id = $1
	`

	e := &datastruct.Netlog{}

	err := r.conn.QueryRow(ctx, query, id).Scan(
		&e.ID,
		&e.CreatedAt,
		&e.Keywords,
		&e.Parameters,
		&e.Error,
		&e.NumBeforeAiFilter,
		&e.NumAfterAiFilter,
		&e.ResultBeforeAiFilter,
		&e.Result,
		&e.ClientID,
		&e.AppName,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // либо можно вернуть кастомную ErrNotFound
		}
		return nil, fmt.Errorf("GetByID scan failed: %w", err)
	}

	return e, nil
}
