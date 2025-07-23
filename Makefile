.PHONY:
dev-tools:
	@mkdir -p $(PWD)/bin
	@curl -sSL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s v2.1.5

infra-tools:
	@mkdir -p $(PWD)/bin
	GOBIN=$(PWD)/bin go install github.com/pressly/goose/v3/cmd/goose@v3.24.2

dev:
	@echo "[RUN] Running the application in development mode..."
	@echo "[RUN] Installing dependencies..."
	@docker compose up --build
	@echo "[RUN] Installing air..."
	@echo "[RUN] Running the application..."

lint-check:
	@echo "[RUN] Running linter checks..."
	@bin/golangci-lint run

lint-fix:
	@echo "[RUN] RRunning linterunning linter fixes..."
	@bin/golangci-lint run --fix

migrate:
	@echo "[RUN] Running database migrations..."
	@sudo go generate ./ent

build:
	@echo "[RUN] Building the application..."
	@echo "[RUN] Installing dependencies..."
	@go mod tidy
	@go build -o build/shinplay ./cmd/server
	@echo "[RUN] Building the application..."

server:
	@echo "[RUN] Running the application in production mode..."
	@echo "[RUN] Building the application..."
	@go build -o build/shinplay ./cmd/server
	@echo "[RUN] Running the application..."
	@./build/shinplay
