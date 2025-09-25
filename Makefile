.PHONY: test test-verbose mock-generate mock-install test-brew clean

# Установка инструментов для генерации моков
mock-install:
	go install github.com/golang/mock/mockgen@latest

hello:
	@echo "Rinat"
	@echo "America"

# Генерация моков
mock-generate:
	@echo "Генерация моков для интерфейсов репозитория..."
	@mkdir -p src/infrastructure/repository/mocks
	mockgen -source=src/infrastructure/repository/interfaces.go -destination=src/infrastructure/repository/mocks/mock_repository.go
	@echo "Моки успешно сгенерированы!"

# Генерация всех моков через go generate
generate:
	go generate ./...

# Запуск всех тестов
test:
	go test ./...

# Запуск тестов с подробным выводом
test-verbose:
	go test -v ./...

# Запуск тестов только для процессора brew
test-brew:
	go test -v ./src/application/processor/brew/

# Запуск тестов с покрытием
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Отчёт о покрытии сохранён в coverage.html"

# Очистка файлов покрытия
clean:
	rm -f coverage.out coverage.html

# Установка зависимостей для тестирования
deps-test:
	go get github.com/golang/mock/gomock@latest
	go get github.com/stretchr/testify/assert@latest

# Полная установка и настройка для разработки
setup: deps-test mock-install mock-generate
	@echo "Окружение для разработки настроено!"

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  setup         - Полная настройка окружения для разработки"
	@echo "  mock-install  - Установка mockgen"
	@echo "  mock-generate - Генерация моков"
	@echo "  generate      - Запуск go generate"
	@echo "  test          - Запуск всех тестов"
	@echo "  test-verbose  - Запуск тестов с подробным выводом"
	@echo "  test-brew     - Запуск тестов процессора brew"
	@echo "  test-coverage - Запуск тестов с отчётом о покрытии"
	@echo "  clean         - Очистка временных файлов"
	@echo "  deps-test     - Установка зависимостей для тестирования"
	@echo "  help          - Показать эту помощь" 