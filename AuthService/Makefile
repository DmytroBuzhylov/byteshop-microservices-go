BINARY_NAME=auth_service
PROTO_SRC=proto
OUT_DIR=generated

PORT=50051
DB_URL=postgres://postgres:root@localhost:5432/BiteShop

build:
	@echo "Компиляция проекта"
	go build -o $(BINARY_NAME) ./cmd/main.go

run: build
	@echo "Запуск сервиса"
	./$(BINARY_NAME)

test:
	@echo "Запуск тестов"
	go test ./...

clean:
	@echo "Очистка проекта"
	rm -f $(BINARY_NAME)

lint:
	@echo "Проверка кода линтером"
	golangci-lint run

docker-build:
	@echo "Сборка Docker-образа"
	docker build -t auth-service .

docker-run:
	@echo "Запуск контейнера"
	docker run -p $(PORT):$(PORT) -e DB_URL=$(DB_URL) auth-service

generate:
	@echo "🔧 Генерация gRPC-кода"
	mkdir -p $(OUT_DIR)/auth $(OUT_DIR)/product $(OUT_DIR)/order $(OUT_DIR)/payment $(OUT_DIR)/notification
	protoc --go_out=$(OUT_DIR)/auth --go-grpc_out=$(OUT_DIR)/auth $(PROTO_SRC)/auth.proto
	protoc --go_out=$(OUT_DIR)/product --go-grpc_out=$(OUT_DIR)/product $(PROTO_SRC)/product.proto
	protoc --go_out=$(OUT_DIR)/order --go-grpc_out=$(OUT_DIR)/order $(PROTO_SRC)/order.proto
	protoc --go_out=$(OUT_DIR)/payment --go-grpc_out=$(OUT_DIR)/payment $(PROTO_SRC)/payment.proto
	protoc --go_out=$(OUT_DIR)/notification --go-grpc_out=$(OUT_DIR)/notification $(PROTO_SRC)/notification.proto




help:
	@echo "Доступные команды:"
	@echo "  make build        - Скомпилировать проект"
	@echo "  make run          - Запустить сервис"
	@echo "  make test         - Запустить тесты"
	@echo "  make clean        - Очистить проект"
	@echo "  make lint         - Проверка кода линтером"
	@echo "  make docker-build - Собрать Docker-образ"
	@echo "  make docker-run   - Запустить контейнер"
