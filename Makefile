DSN=root:root@tcp(localhost:3306)/employee_db?parseTime=true

.PHONY: dev dev-local infra down setup-mac

setup-mac:
	@which brew > /dev/null || (echo "Instalando Homebrew..." && /bin/bash -c "$$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)")
	@which go > /dev/null || (echo "Instalando Go..." && brew install go)
	@which temporal > /dev/null || (echo "Instalando Temporal CLI..." && brew install temporal)
	@echo "Tudo pronto! Rode: make dev-local"

infra:
	docker compose up -d

down:
	docker compose down

dev-local:
	-pkill -f "temporal server" 2>/dev/null; true
	-pkill -f "cmd/api" 2>/dev/null; true
	temporal server start-dev &
	sleep 3
	DATABASE_DSN="$(DSN)" go run ./cmd/api & \
	DATABASE_DSN="$(DSN)" go run ./cmd/worker

dev: infra
	-fuser -k 8080/tcp 2>/dev/null; true
	DATABASE_DSN="$(DSN)" go run ./cmd/api & \
	DATABASE_DSN="$(DSN)" go run ./cmd/worker
