package datastruct

import (
	"time"
)

// Netlog структура для таблицы
type Netlog struct {
	ID                   int64                  `db:"id"`
	CreatedAt            time.Time              `db:"created_at"`
	Keywords             []string               `db:"keywords"`
	Parameters           map[string]interface{} `db:"parameters"`
	Error                *string                `db:"error"`
	NumBeforeAiFilter    int                    `db:"num_before_ai_filter"`
	NumAfterAiFilter     int                    `db:"num_after_ai_filter"`
	ResultBeforeAiFilter map[string]interface{} `db:"result_before_ai_filter"`
	Result               map[string]interface{} `db:"result"`
	ClientID             string                 `db:"client_id"`
	AppName              string                 `db:"app_name"`
}

// Структура фильтра (можно расширять по потребности)
type NetlogFilter struct {
	ClientID    string     `json:"client_id,omitempty"`
	AppName     *string    `json:"app_name,omitempty"`
	Keywords    []string   `json:"keywords,omitempty"`  // точное совпадение всех
	HasError    *bool      `json:"has_error,omitempty"` // true = error IS NOT NULL
	MinBeforeAI *int       `json:"min_before_ai,omitempty"`
	MaxBeforeAI *int       `json:"max_before_ai,omitempty"`
	FromTime    *time.Time `json:"from_time,omitempty"`
	ToTime      *time.Time `json:"to_time,omitempty"`
	Limit       int        `json:"limit"`             // обязательно
	LastID      *int64     `json:"last_id,omitempty"` // для пагинации "после этого id"
}
