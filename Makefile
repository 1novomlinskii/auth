.PHONY: run build test lint make-local-up make-local-down make-local-clear

run:
	go run ./cmd/auth/

build:
	go build -ldflags="-s -w" -o bin/auth ./cmd/auth/

# Произвести генерацию
# Проверка линтера. F=1 флаг автофикса
.PHONY: lint
lint:
	go vet ./...
	golangci-lint run ./... $(if $(F),--fix,)

# Запускает все тесты
.PHONY: test
test:
	go test -v -race -count=1 ./...

# Usage:
#   make-local-up            — docker compose up -d
#   make-local-up 1          — docker compose up -d --build
make-local-up:
	@if [ "$(filter-out $@,$(MAKECMDGOALS))" = "1" ]; then \
		docker compose up -d --build; \
	else \
		docker compose up -d; \
	fi

make-local-down:
	docker compose down

make-local-clear:
	docker compose down -v
