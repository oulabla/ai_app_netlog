# =============================================================================
# Makefile
# =============================================================================
-include .project.mk

# =============================================================================
# Загрузка .env (если файл существует)
# =============================================================================
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# =============================================================================
# Основные переменные
# =============================================================================
PROJECT_NAME ?= go-base
MODULE       ?= github.com/oulabla/$(PROJECT_NAME)
BIN_DIR      := $(CURDIR)/bin
PROTO_DIR    := proto
GEN_DIR      := gen
CURRENT_DATE := $(shell date +%Y-%m-%d)

# =============================================================================
# Инструменты и версии
# =============================================================================
GOOSE_VERSION          := v3.27.0
BUF_VERSION            := v1.66.0          # ← можно обновить до v1.66.0 при желании
YQ_VERSION             := v4.52.4
PROTOC_GEN_GO          := v1.36.0
PROTOC_GEN_GO_GRPC     := v1.5.1
GRPC_GATEWAY_VERSION   := v2.23.0
OPENAPI_VERSION        := v2.20.0
MOCKERY_VERSION        := v3.6.4
GOLANGCI_LINT_VERSION  := v2.10.1
GOLANGCI_LINT          := $(BIN_DIR)/golangci-lint

# =============================================================================
# Пути к бинарникам
# =============================================================================
BUF      := $(BIN_DIR)/buf
YQ       := $(BIN_DIR)/yq
MOCKERY  := $(BIN_DIR)/mockery
GOOSE    := $(BIN_DIR)/goose

# =============================================================================
.DEFAULT_GOAL := help
.PHONY: help init tools buf-update generate generate-config proto-generate scaffold clean test lint tidy add-service build run \
        migrate migrate-new migrate-up migrate-down migrate-status migrate-reset goose

help:
	@echo ""
	@echo "Доступные команды:"
	@echo ""
	@echo "  make init              Инициализация проекта (tools + buf-update + tidy + config + generate)"
	@echo "  make tools             Установить все инструменты в ./bin"
	@echo "  make buf-update        Обновить buf-зависимости (deps)"
	@echo "  make generate          Генерация proto + stubs + config keys (без установки инструментов и без buf dep update)"
	@echo "  make generate-config   Только config keys"
	@echo "  make proto-generate    Только protobuf + grpc-gateway + openapi (без buf dep update)"
	@echo "  make build             Сборка → ./bin/$(PROJECT_NAME)"
	@echo "  make run               Запуск"
	@echo "  make test              Тесты + покрытие"
	@echo "  make lint              golangci-lint"
	@echo "  make tidy              go mod tidy"
	@echo "  make clean             Очистка"
	@echo ""
	@echo "  make add-service users      → proto/users/v1/..."
	@echo "  make add-service users v2   → proto/users/v2/..."
	@echo ""
	@echo "Миграции (goose):"
	@echo "  make migrate new      <имя_миграции>     → создать новую миграцию"
	@echo "  make migrate up                          → применить все ожидающие"
	@echo "  make migrate down                        → откатить последнюю"
	@echo "  make migrate status                      → показать статус миграций"
	@echo "  make migrate reset                       → откатить ВСЕ миграции (с подтверждением)"
	@echo ""

# =============================================================================
# Инициализация проекта
# =============================================================================
init: tools buf-update tidy generate-config generate
	@echo ""
	@echo "Инициализация завершена ✓"
	@echo "Можно запускать: make run"

# =============================================================================
# Tools
# =============================================================================
tools: $(BIN_DIR) buf yq protoc-plugins mockery golangci-lint goose

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

buf: $(BIN_DIR)
	curl -sSL \
	  "https://github.com/bufbuild/buf/releases/latest/download/buf-Linux-x86_64" \
	  -o $(BUF)
	chmod +x $(BUF)
	@echo "buf latest установлен → $$( $(BUF) --version )"

yq: $(BIN_DIR)
	curl -sSL \
	  "https://github.com/mikefarah/yq/releases/download/$(YQ_VERSION)/yq_linux_amd64" \
	  -o $(YQ)
	chmod +x $(YQ)
	@echo "yq v$(YQ_VERSION) → $(YQ)"

protoc-plugins: $(BIN_DIR)
	GOBIN=$(BIN_DIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO)
	GOBIN=$(BIN_DIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC)
	GOBIN=$(BIN_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@$(GRPC_GATEWAY_VERSION)
	GOBIN=$(BIN_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@$(OPENAPI_VERSION)

mockery: $(BIN_DIR)
	GOBIN=$(BIN_DIR) go install github.com/vektra/mockery/v3@$(MOCKERY_VERSION)

golangci-lint: $(BIN_DIR)
	@echo "→ Installing golangci-lint v$(GOLANGCI_LINT_VERSION)"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
	  | sh -s -- -b $(BIN_DIR) v$(GOLANGCI_LINT_VERSION)
	@$(GOLANGCI_LINT) --version

goose: $(BIN_DIR)
	@echo "→ Installing goose $(GOOSE_VERSION)"
	GOBIN=$(BIN_DIR) go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)
	@$(GOOSE) version

# =============================================================================
# Buf зависимости — отдельная цель, не вызывается автоматически
# =============================================================================
buf-update:
	@if [ ! -x "$(BUF)" ]; then \
		echo "→ buf не найден. Запустите: make tools"; \
		exit 1; \
	fi
	@echo "→ Updating buf dependencies..."
	PATH="$(CURDIR)/bin:$$PATH" $(BUF) dep update $(PROTO_DIR)

# =============================================================================
# Генерация (теперь без автоматической установки и buf-update)
# =============================================================================
generate: proto-generate scaffold generate-config
	@echo "→ Генерация завершена"

generate-config:
	@echo "→ Generating config keys..."
	@bash script/generate_config.sh

proto-generate:
	@if [ ! -x "$(BUF)" ]; then \
		echo "→ buf не найден. Запустите: make tools"; \
		exit 1; \
	fi
	@echo "→ Generating protobuf files..."
	PATH="$(CURDIR)/bin:$$PATH" $(BUF) generate $(PROTO_DIR)

scaffold:
	@PROJECT_NAME="$(PROJECT_NAME)" MODULE="$(MODULE)" \
		bash script/gen_stub.sh

# =============================================================================
# Остальные цели
# =============================================================================
tidy:
	go mod tidy

test: tidy
	@echo "→ Запуск тестов с покрытием..."
	@if [ -z "$(COVER_PKGS)" ]; then \
		go test -v -race ./...; \
		exit 0; \
	fi
	go test -v -race \
		-coverpkg=$(COVER_PKGS) \
		-coverprofile=coverage.out \
		-covermode=atomic \
		./...
	@go tool cover -func=coverage.out | grep total || echo "→ Покрытие 0.0%"
	@go tool cover -html=coverage.out -o coverage.html 2>/dev/null || true

lint:
	@if [ ! -x "$(GOLANGCI_LINT)" ]; then \
		echo "→ golangci-lint не найден. Запусти: make tools"; \
		exit 1; \
	fi
	PATH="$(BIN_DIR):$$PATH" $(GOLANGCI_LINT) run --timeout=5m --color=always

build: tidy generate
	CGO_ENABLED=0 go build \
		-trimpath \
		-ldflags "-s -w -X main.version=$(CURRENT_DATE)" \
		-o $(BIN_DIR)/$(PROJECT_NAME) \
		./cmd/$(PROJECT_NAME)/main.go
	@echo "→ Собрано: $(BIN_DIR)/$(PROJECT_NAME)"

run: build
	@$(BIN_DIR)/$(PROJECT_NAME)

# =============================================================================
# Добавление сервиса
# =============================================================================
.PHONY: add-service
add-service: ; @bash script/add_service.sh $(filter-out $@,$(MAKECMDGOALS))
%:
	@:

# =============================================================================
# Очистка
# =============================================================================
clean:
	rm -rf $(GEN_DIR) $(BIN_DIR) coverage.out coverage.html

# =============================================================================
# Миграции (Goose + DB_URL из .env)
# =============================================================================
migrate:
	@if [ -z "$$DB_URL" ]; then \
		echo "Ошибка: переменная DB_URL не задана"; \
		echo "Добавьте в .env строку вида:"; \
		echo "  DB_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable"; \
		echo ""; \
		echo "Доступные команды миграций:"; \
		echo "  make migrate new      <имя>"; \
		echo "  make migrate up"; \
		echo "  make migrate down"; \
		echo "  make migrate status"; \
		echo "  make migrate reset"; \
		exit 1; \
	fi
	@:

migrate-new: migrate
	@name="$(filter-out $@,$(MAKECMDGOALS))"; \
	if [ -z "$$name" ]; then \
		echo "Ошибка: укажите имя миграции"; \
		echo "Пример:  make migrate new create-users-table"; \
		exit 1; \
	fi; \
	$(GOOSE) -dir migrations create "$$name" sql

migrate-up: migrate
	@if [ ! -x "$(GOOSE)" ]; then \
		echo "→ goose не найден. Запустите: make tools"; \
		exit 1; \
	fi
	@echo "→ goose up..."
	@$(GOOSE) -dir migrations postgres "$$DB_URL" up

migrate-down: migrate
	@if [ ! -x "$(GOOSE)" ]; then \
		echo "→ goose не найден. Запустите: make tools"; \
		exit 1; \
	fi
	@echo "→ goose down..."
	@$(GOOSE) -dir migrations postgres "$$DB_URL" down

migrate-status: migrate
	@if [ ! -x "$(GOOSE)" ]; then \
		echo "→ goose не найден. Запустите: make tools"; \
		exit 1; \
	fi
	@echo "→ goose status:"
	@$(GOOSE) -dir migrations postgres "$$DB_URL" status

migrate-reset: migrate
	@if [ ! -x "$(GOOSE)" ]; then \
		echo "→ goose не найден. Запустите: make tools"; \
		exit 1; \
	fi
	@echo "ВНИМАНИЕ! Будут откатаны ВСЕ миграции и потеряны данные."
	@read -p "Подтвердите (y/N): " confirm && [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ] || (echo "Отменено"; exit 1)
	@$(GOOSE) -dir migrations postgres "$$DB_URL" reset