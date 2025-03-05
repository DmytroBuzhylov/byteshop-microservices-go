BINARY_NAME=product_service

PORT=50052
DB_URL=

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
	docker build -t product-service .

docker-run:
	@echo "Запуск контейнера"
	docker run -p $(PORT):$(PORT) -e DB_URL=$(DB_URL) product-service


help:
	@echo "Доступные команды:"
	@echo "  make build        - Скомпилировать проект"
	@echo "  make run          - Запустить сервис"
	@echo "  make test         - Запустить тесты"
	@echo "  make clean        - Очистить проект"
	@echo "  make lint         - Проверка кода линтером"
	@echo "  make docker-build - Собрать Docker-образ"
	@echo "  make docker-run   - Запустить контейнер"
