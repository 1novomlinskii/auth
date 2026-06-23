
.PHONY: run
run:
	go run ./cmd/auth/

.PHONY: build
build:
	go build -ldflags="-s -w" -o bin/auth ./cmd/auth/

# Генерация proto + qtc
.PHONY: gen
gen:
	buf generate
	qtc -dir=internal/repository/pg/query

# Проверка линтера. F=1 флаг автофикса
.PHONY: lint
lint:
	go vet ./...
	golangci-lint run ./... $(if $(F),--fix,)

# Запускает все тесты
.PHONY: test
test:
	go test -v -race -count=1 ./...

# Поднимает окружение для локальной разработки (make local-up B=1 для build флага)
.PHONY: local-up
local-up:
	 $(COMPOSE_CMD) up $(if $(B),--build,)

.PHONY: local-down
local-down:
	docker compose down

.PHONY: local-clear
local-clear:
	docker compose down -v
	
