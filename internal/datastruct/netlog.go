package datastruct

import "time"

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
