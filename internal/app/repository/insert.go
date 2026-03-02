package repository

import (
	"context"
	"fmt"

	"github.com/oulabla/ai_app_netlog/internal/datastruct"
)

func (r *Repository) Insert(ctx context.Context, netlog *datastruct.Netlog) (int64, error) {
	query := `
        INSERT INTO netlog (
            created_at, keywords, parameters, error,
            num_before_ai_filter, num_after_ai_filter,
            result_before_ai_filter, result, client_id, app_name
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id
    `

	var id int64
	err := r.conn.QueryRow(ctx, query,
		netlog.CreatedAt,
		netlog.Keywords,
		netlog.Parameters,
		netlog.Error,
		netlog.NumBeforeAiFilter,
		netlog.NumAfterAiFilter,
		netlog.ResultBeforeAiFilter,
		netlog.Result,
		netlog.ClientID,
		netlog.AppName,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка вставки: %w", err)
	}

	return id, nil
}
