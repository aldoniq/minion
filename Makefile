.PHONY: build clean test run-extend run-refresh install help fmt vet

# Переменные
BINARY_NAME=minion
BUILD_DIR=build
CMD_DIR=cmd/minion

# По умолчанию - показать справку
help:
	@echo "🍌 Minion - Команды сборки:"
	@echo "  build        - Собрать исполняемый файл"
	@echo "  install      - Установить в GOPATH/bin"
	@echo "  clean        - Удалить скомпилированные файлы"
	@echo "  test         - Запустить тесты"
	@echo "  fmt          - Отформатировать код"
	@echo "  vet          - Проверить код с go vet"
	@echo "  run-extend   - Запустить продление API ключей"
	@echo "  run-refresh  - Запустить обновление меню"
	@echo "  deps         - Установить зависимости"

# Сборка
build:
	@echo "🔨 Сборка minion..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "✅ Готово: $(BUILD_DIR)/$(BINARY_NAME)"

# Установка
install:
	@echo "📦 Установка minion..."
	@go install ./$(CMD_DIR)
	@echo "✅ Установлено в GOPATH/bin"

# Очистка
clean:
	@echo "🧹 Очистка..."
	@rm -rf $(BUILD_DIR)
	@echo "✅ Очищено"

# Тесты
test:
	@echo "🧪 Запуск тестов..."
	@go test -v ./...

# Форматирование кода
fmt:
	@echo "💄 Форматирование кода..."
	@go fmt ./...

# Проверка кода
vet:
	@echo "🔍 Проверка кода..."
	@go vet ./...

# Установка зависимостей
deps:
	@echo "📥 Установка зависимостей..."
	@go mod tidy
	@go mod download

# Запуск команд
run-extend:
	@echo "🔑 Запуск продления API ключей..."
	@go run ./$(CMD_DIR) extend-keys

run-refresh:
	@echo "🍽️ Запуск обновления меню..."
	@go run ./$(CMD_DIR) refresh-menus

# Разработка
dev-build: fmt vet build

# Полная проверка
check: fmt vet test 