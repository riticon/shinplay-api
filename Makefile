.PHONY:
tools:
	@curl -sSL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s

dev:
	@echo "[RUN] Running the application in development mode..."
	@echo "[RUN] Installing dependencies..."
	@go mod tidy
	@echo "[RUN] Installing air..."
	@echo "[RUN] Running the application..."
	@bin/air

build:
	@echo "[RUN] Building the application..."
	@echo "[RUN] Installing dependencies..."
	@go mod tidy
	@go build -o build/shinplay ./cmd/api
	@echo "[RUN] Building the application..."
