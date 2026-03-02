-- +goose Up
-- +goose StatementBegin
CREATE TABLE netlog (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    keywords TEXT[],
    parameters JSONB,
    error TEXT,
    num_before_ai_filter INTEGER NOT NULL DEFAULT 0,
    num_after_ai_filter INTEGER NOT NULL DEFAULT 0,
    result_before_ai_filter JSONB,
    result JSONB,
    client_id UUID NOT NULL,
    app_name TEXT NOT NULL
);

-- Комментарии
COMMENT ON TABLE netlog IS 'Таблица для логирования AI запросов и обработки';
COMMENT ON COLUMN netlog.id IS 'Уникальный идентификатор записи';
COMMENT ON COLUMN netlog.created_at IS 'Дата и время создания записи с часовым поясом';
COMMENT ON COLUMN netlog.keywords IS 'Массив ключевых слов для поиска и анализа';
COMMENT ON COLUMN netlog.parameters IS 'Параметры запроса в формате JSONB';
COMMENT ON COLUMN netlog.error IS 'Текст ошибки, если произошла';
COMMENT ON COLUMN netlog.num_before_ai_filter IS 'Количество элементов до AI фильтрации';
COMMENT ON COLUMN netlog.num_after_ai_filter IS 'Количество элементов после AI фильтрации';
COMMENT ON COLUMN netlog.result_before_ai_filter IS 'Результат до применения AI фильтров (JSONB)';
COMMENT ON COLUMN netlog.result IS 'Итоговый результат после обработки (JSONB)';
COMMENT ON COLUMN netlog.client_id IS 'ID клиента в формате UUID';
COMMENT ON COLUMN netlog.app_name IS 'Название приложения';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS netlog;
-- +goose StatementEnd