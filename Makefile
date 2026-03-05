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
BUF_VERSION            := v1.66.0
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
        migrate migrate-new migrate-up migrate-down migrate-status migrate-reset goose \
        mocks mocks-quick update-mockery-config generate-mocks list-mocks clean-mocks

help:
	@echo ""
	@echo "Доступные команды:"
	@echo ""
	@echo "Основные:"
	@echo "  make init              → Инициализация (tools + buf-update + tidy + config + generate)"
	@echo "  make tools             → Установить все инструменты в ./bin"
	@echo "  make generate          → Генерация proto + stubs + config keys"
	@echo "  make build             → Сборка → ./bin/$(PROJECT_NAME)"
	@echo "  make run               → Запуск"
	@echo "  make test              → Тесты + покрытие"
	@echo "  make lint              → golangci-lint"
	@echo "  make tidy              → go mod tidy"
	@echo "  make clean             → Очистка"
	@echo ""
	@echo "Сервисы:"
	@echo "  make add-service users      → Создать прото-файл с заглушкой proto/users/v1/..."
	@echo "  make add-service users v2   → Создать прото-файл с заглушкой proto/users/v2/..."
	@echo ""
	@echo "Моки (mockery):"
	@echo "  make mocks                  → Обновить .mockery.yaml + сгенерировать все моки (рекомендуется)"
	@echo "  make mocks-quick            → Сгенерировать моки без обновления конфига (быстрее)"
	@echo "  make update-mockery-config  → Только обновить .mockery.yaml"
	@echo "  make list-mocks             → Показать найденные интерфейсы и их мок-файлы"
	@echo "  make clean-mocks            → Удалить все *_mock.go файлы"
	@echo ""
	@echo "Миграции (goose):"
	@echo "  make migrate new      <имя_миграции>     → создать новую"
	@echo "  make migrate up                          → применить ожидающие"
	@echo "  make migrate down                        → откатить последнюю"
	@echo "  make migrate status                      → статус"
	@echo "  make migrate reset                       → откатить ВСЕ (с подтверждением)"
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
	@echo "buf latest → $$( $(BUF) --version )"

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
# Buf зависимости
# =============================================================================
buf-update:
	@if [ ! -x "$(BUF)" ]; then \
		echo "→ buf не найден. Запустите: make tools"; \
		exit 1; \
	fi
	@echo "→ Updating buf dependencies..."
	PATH="$(CURDIR)/bin:$$PATH" $(BUF) dep update $(PROTO_DIR)

# =============================================================================
# Генерация
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
# Mocks (mockery)
# =============================================================================
MOCKERY_CONFIG := .mockery.yaml

update-mockery-config:
	@echo "→ Updating mockery configuration (.mockery.yaml)..."
	@bash script/update-mockery-config.sh
	@echo "→ Конфиг обновлён"

mocks: update-mockery-config
	@echo "→ Генерация моков (с актуальным конфигом)..."
	@if [ ! -x "$(MOCKERY)" ]; then \
		echo "❌ mockery не найден → запустите: make mockery"; \
		exit 1; \
	fi
	@if [ ! -f "$(MOCKERY_CONFIG)" ]; then \
		echo "❌ $(MOCKERY_CONFIG) не найден → сначала: make update-mockery-config"; \
		exit 1; \
	fi
	@$(MOCKERY) --config $(MOCKERY_CONFIG)
	@echo "→ Моки сгенерированы ✓"

mocks-quick:
	@echo "→ Быстрая генерация моков (без обновления конфига)..."
	@if [ ! -x "$(MOCKERY)" ]; then \
		echo "❌ mockery не найден → запустите: make mockery"; \
		exit 1; \
	fi
	@if [ ! -f "$(MOCKERY_CONFIG)" ]; then \
		echo "❌ $(MOCKERY_CONFIG) не найден → сначала: make update-mockery-config"; \
		exit 1; \
	fi
	@$(MOCKERY) --config $(MOCKERY_CONFIG)
	@echo "→ Быстрая генерация завершена"

list-mocks:
	@echo "Найденные *Mock интерфейсы и их мок-файлы:"
	@find ./internal -type f -name "*Mock.go" -not -path "*/mocks/*" | while read -r iface_file; do \
		iface=$$(grep -E "^type\s+[A-Za-z0-9_]+Mock\s+interface" "$$iface_file" | sed -E 's/^type\s+([A-Za-z0-9_]+)Mock.*/\1/'); \
		if [ -n "$$iface" ]; then \
			mock_file="$$(dirname "$$iface_file")/mocks/$${iface}_mock.go"; \
			if [ -f "$$mock_file" ]; then \
				echo "  ✔  $$iface\t→ $$mock_file"; \
			else \
				echo "  ✗  $$iface\t(мок ещё не сгенерирован)"; \
			fi; \
		fi; \
	done
	@echo ""
	@echo "Всего мок-файлов: $$(find ./internal -type f -name "*_mock.go" | wc -l)"

clean-mocks:
	@echo "→ Удаление всех сгенерированных *_mock.go файлов..."
	@find ./internal -type f -name "*_mock.go" -delete
	@echo "→ Удалено $$(find ./internal -type f -name "*_mock.go" 2>/dev/null | wc -l) файлов"
	@echo "→ Рекомендуется: make update-mockery-config && make mocks"

# =============================================================================
# Остальные цели
# =============================================================================
tidy:
	go mod tidy

test: tidy
	@echo "→ Тесты с покрытием"
	go test -v -race \
		-coverpkg=$(shell go list ./internal/... | grep -vE '/(metric|server|script)$$' | tr '\n' ',' | sed 's/,$$//') \
		-coverprofile=coverage.out \
		-covermode=atomic \
		./...
	@go tool cover -func=coverage.out | grep total

lint:
	@if [ ! -x "$(GOLANGCI_LINT)" ]; then \
		echo "→ golangci-lint не найден. Запустите: make tools"; \
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