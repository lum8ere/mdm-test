COMPOSE_FILE=docker-compose.yml

up:
	docker-compose -f $(COMPOSE_FILE) up -d

down:
	docker-compose -f $(COMPOSE_FILE) down

restart:
	@make down
	@make up

# ТЕСТЫ ДЛЯ БЕКА
# test-backend:
# 	cd backend && go test -v ./...

# Запуск серверной части (бэкенда)
run-backend:
	cd backend/app/backend-api && go run backend-api.go

# Запуск клиентской части (агента)
run-client:
	cd client && go run main.go --device-id=android-test --server=http://localhost:4000