DSN=root:root@tcp(localhost:3306)/employee_db?parseTime=true

.PHONY: dev infra down

infra:
	docker compose up -d

down:
	docker compose down

dev: infra
	-fuser -k 8080/tcp 2>/dev/null; true
	DATABASE_DSN="$(DSN)" go run ./cmd/api & \
	DATABASE_DSN="$(DSN)" go run ./cmd/worker
