# Переменные
BINARY_NAME=minion
BUILD_DIR=bin
CMD_DIR=cmd/minion

# По умолчанию - показать справку
help:
	@echo "🍌 Minion - Команды сборки:"
	@echo "  build        - Собрать исполняемый файл"
	@echo "  run          - Запустить приложение"
	@echo "  fmt          - Отформатировать код"
	@echo "  vet          - Проверить код с go vet"
	@echo ""
	@echo "🔑 Команды приложения:"
	@echo "  extend-keys  - Продление API ключей"
	@echo "  refresh-menus - Обновление меню"
	@echo ""
	@echo "🔧 Команды настройки:"
	@echo "  deps         - Установить зависимости"

# Сборка
build:
	@echo "🔨 Сборка minion..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "✅ Готово: $(BUILD_DIR)/$(BINARY_NAME)"


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

# Продление API ключей
extend-keys:
	@echo "🔑 Продление API ключей..."
	@./$(BUILD_DIR)/$(BINARY_NAME) extend-keys || go run ./$(CMD_DIR) extend-keys

# Обновление меню
refresh-menus:
	@echo "🍽️ Обновление меню..."
	@./$(BUILD_DIR)/$(BINARY_NAME) refresh-menus || go run ./$(CMD_DIR) refresh-menus