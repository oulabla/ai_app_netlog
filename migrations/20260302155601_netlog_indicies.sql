-- +goose Up
-- +goose NO TRANSACTION
-- Индексы
CREATE INDEX concurrently idx_netlog_created_at ON netlog(created_at);
CREATE INDEX concurrently idx_netlog_client_id ON netlog(client_id);
CREATE INDEX concurrently idx_netlog_app_name ON netlog(app_name);

-- GIN индексы
CREATE INDEX concurrently idx_netlog_keywords ON netlog USING GIN (keywords);

-- Составные индексы
CREATE INDEX concurrently idx_netlog_client_created ON netlog(client_id, created_at);
CREATE INDEX concurrently idx_netlog_num_before ON netlog(num_before_ai_filter);
CREATE INDEX concurrently idx_netlog_num_after ON netlog(num_after_ai_filter);

-- +goose Down
-- +goose NO TRANSACTION
drop index concurrently if exists idx_netlog_created_at;
drop index concurrently if exists idx_netlog_client_id;
drop index concurrently if exists idx_netlog_app_name;

drop index concurrently if exists idx_netlog_keywords;

drop index concurrently if exists idx_netlog_client_created;
drop index concurrently if exists idx_netlog_num_before;
drop index concurrently if exists idx_netlog_num_after;