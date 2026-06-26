DSN=employee:employee@tcp(localhost:3306)/employee_db?parseTime=true

.PHONY: dev dev-local infra down setup-mac seed

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
	temporal server start-dev & \
	echo "Aguardando Temporal..." && \
	until nc -z localhost 7233 2>/dev/null; do sleep 1; done && \
	echo "Temporal pronto!" && \
	DATABASE_DSN="$(DSN)" go run ./cmd/api & \
	DATABASE_DSN="$(DSN)" go run ./cmd/worker

seed:
	@echo "Inserindo 400 registros..."
	@bash -c '\
	NAMES=(Ana Bruno Carla Diego Elena Felipe Gabriela Hugo Isabela Joao Karen Lucas Mariana Nicolas Olivia Pedro Rafael Sofia Thiago Camila); \
	SURNAMES=(Silva Santos Oliveira Souza Lima Pereira Costa Ferreira Alves Rodrigues Martins Gomes Barbosa Carvalho Rocha Dias Nascimento Andrade Moreira Nunes); \
	POSITIONS=(Engineer Designer Manager Analyst DevOps QA Architect); \
	TYPES=(fulltime contractor); \
	for i in $$(seq 1 400); do \
		NAME="$${NAMES[$$((RANDOM % $${#NAMES[@]}))]}_$${SURNAMES[$$((RANDOM % $${#SURNAMES[@]}))]}"; \
		NAME=$${NAME//_/ }; \
		POSITION=$${POSITIONS[$$((RANDOM % $${#POSITIONS[@]}))]}; \
		TYPE=$${TYPES[$$((RANDOM % $${#TYPES[@]}))]}; \
		SALARY=$$((3000 + RANDOM % 12000)); \
		curl -s -X POST http://localhost:8080/employees \
			-H "Content-Type: application/json" \
			-d "{\"name\":\"$$NAME\",\"position\":\"$$POSITION\",\"salary\":$$SALARY,\"type\":\"$$TYPE\"}" > /dev/null; \
	done'
	@echo "400 registros inseridos!"

dev: infra
	-fuser -k 8080/tcp 2>/dev/null; true
	DATABASE_DSN="$(DSN)" go run ./cmd/api & \
	DATABASE_DSN="$(DSN)" go run ./cmd/worker
