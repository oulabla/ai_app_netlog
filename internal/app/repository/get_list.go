package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/oulabla/ai_app_netlog/internal/datastruct"
)

// возвращаем записи + последний id + ошибка
func (r *Repository) GetList(ctx context.Context, filter *datastruct.NetlogFilter) ([]*datastruct.Netlog, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Limit > 500 {
		filter.Limit = 500 // разумный лимит
	}

	var rows pgx.Rows
	var err error

	// Базовая часть запроса
	query := `
		SELECT 
			id,
			created_at,
			keywords,
			parameters,
			error,
			num_before_ai_filter,
			num_after_ai_filter,
			client_id,
			app_name
		FROM netlog
		WHERE 1=1
	`

	args := []any{}

	// 1. Фильтр по last_id (пагинация вперед)
	if filter.LastID != nil && *filter.LastID > 0 {
		query += ` AND id < $1`
		args = append(args, *filter.LastID)
	}

	// 2. Фильтры
	if filter.ClientID != "" {
		query += fmt.Sprintf(" AND client_id = $%d", len(args)+1)
		args = append(args, filter.ClientID)
	}

	if filter.AppName != nil && *filter.AppName != "" {
		query += fmt.Sprintf(" AND app_name = $%d", len(args)+1)
		args = append(args, *filter.AppName)
	}

	if len(filter.Keywords) > 0 {
		// @> — содержит все переданные элементы (массив содержит массив)
		query += fmt.Sprintf(" AND keywords @> $%d", len(args)+1)
		args = append(args, filter.Keywords)
	}

	if filter.HasError != nil {
		if *filter.HasError {
			query += " AND error IS NOT NULL"
		} else {
			query += " AND error IS NULL"
		}
	}

	if filter.MinBeforeAI != nil {
		query += fmt.Sprintf(" AND num_before_ai_filter >= $%d", len(args)+1)
		args = append(args, *filter.MinBeforeAI)
	}

	if filter.MaxBeforeAI != nil {
		query += fmt.Sprintf(" AND num_before_ai_filter <= $%d", len(args)+1)
		args = append(args, *filter.MaxBeforeAI)
	}

	if filter.FromTime != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", len(args)+1)
		args = append(args, filter.FromTime)
	}

	if filter.ToTime != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", len(args)+1)
		args = append(args, filter.ToTime)
	}

	// Сортировка — почти всегда по убыванию времени / id
	query += `
		ORDER BY id DESC
		LIMIT $` + fmt.Sprintf("%d", len(args)+1)

	args = append(args, filter.Limit)

	// ────────────────────────────────────────────────
	// Выполнение
	// ────────────────────────────────────────────────

	rows, err = r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var items []*datastruct.Netlog
	var lastID int64 = 0

	for rows.Next() {
		e := &datastruct.Netlog{}

		err = rows.Scan(
			&e.ID,
			&e.CreatedAt,
			&e.Keywords,
			&e.Parameters,
			&e.Error,
			&e.NumAfterAiFilter,
			&e.NumAfterAiFilter,
			&e.ClientID,
			&e.AppName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan failed: %w", err)
		}

		items = append(items, e)
		lastID = e.ID
	}

	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}

	// Если ничего не нашли — lastID остаётся 0
	return items, lastID, nil
}
