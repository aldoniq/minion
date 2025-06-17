# 🍌 Minion - HTTP API для автоматизации iiko

BELLO! HTTP API сервер для автоматизации работы с iiko системами. Продление ключей и обновление меню через веб-интерфейс.

## 🚀 Быстрый старт

### Требования

- Go 1.20+
- AWS CLI настроен (`aws configure`)
- Доступ к AWS Secrets Manager
- MongoDB база данных с ресторанами

### Установка

```bash
# Клонируем репозиторий
git clone <repository-url>
cd minion

# Устанавливаем зависимости
go mod tidy

# Создаем .env файл
cp env.example .env

# Редактируем .env файл
nano .env
```

### Настройка

Создайте `.env` файл на основе `env.example`:

```bash
HTTP_PORT=3000
AWS_REGION=eu-west-1
AWS_SECRET_NAME=ProdEnvs
```

AWS секрет должен содержать:
```json
{
  "db_url": "mongodb://username:password@host:port/database",
  "db_name": "database_name",
  "db_engine": "mongodb"
}
```

### Использование

```bash
# Сборка проекта
go build -o bin/minion ./cmd/minion

# Запуск HTTP API сервера
./bin/minion
# или напрямую
go run ./cmd/minion
```

**HTTP API Эндпоинты:**

| Метод | URL | Описание |
|-------|-----|----------|
| `GET` | `/api/health` | Проверка состояния API |
| `GET` | `/api/config` | Текущая конфигурация |
| `POST` | `/api/extend-keys` | Продление API ключей |
| `POST` | `/api/refresh-menus` | Обновление меню |

**Примеры запросов:**

```bash
# Проверка состояния
curl http://localhost:3000/api/health

# Продление ключей
curl -X POST http://localhost:3000/api/extend-keys

# Обновление меню
curl -X POST http://localhost:3000/api/refresh-menus
```

**Формат ответа:**

```json
{
  "success": true,
  "message": "🎉 GELATO! Обработано 5 ресторанов",
  "data": {
    "processed_restaurants": 5,
    "successful": 4,
    "failed": 1,
    "duration": "2.5s",
    "details": [
      {
        "name": "Ресторан 1",
        "success": true,
        "updated": 3,
        "message": "Обновлено 3 ключей"
      }
    ]
  }
}
```

## 🔧 Конфигурация

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `HTTP_PORT` | Порт HTTP сервера | `3000` |
| `AWS_REGION` | AWS регион | `eu-west-1` |
| `AWS_SECRET_NAME` | Имя секрета в AWS | `ProdEnvs` |

### Структура базы данных

Приложение ожидает коллекцию `restaurants` в MongoDB со следующей структурой:

```javascript
{
  "_id": ObjectId,
  "name": "Название ресторана",
  "pos_type": "iiko",
  "is_deleted": false,
  "iiko_cloud": {
    "custom_domain": "restaurant.iikoweb.ru",
    "login": "iiko_login",
    "password": "iiko_password",
    "organization_id": "...",
    "terminal_id": "...",
    "key": "..."
  },
  "settings": {
    "is_deleted": false
  }
}
```

## 🏗️ Архитектура (KISS принцип)

```
cmd/minion/           - Точка входа (только HTTP сервер)
internal/
├── aws/             - AWS Secrets Manager клиент
├── client/          - HTTP клиент для iiko API
├── config/          - Конфигурация и загрузка ресторанов
├── database/        - MongoDB сервис
├── handlers/        - HTTP API handlers (Fiber)
├── server/          - HTTP сервер (Fiber)
└── models/          - Структуры данных
```

## 🛠️ Команды разработки

```bash
# Сборка
go build -o bin/minion ./cmd/minion

# Запуск
go run ./cmd/minion

# Тесты
go test ./...

# Форматирование
go fmt ./...

# Проверка кода
go vet ./...

# Установка зависимостей
go mod tidy
```

## 🔒 Безопасность

- 🔐 AWS Secrets Manager для безопасного хранения credentials
- 🚫 `.env` файлы добавлены в `.gitignore`
- 👥 Индивидуальные iiko credentials для каждого ресторана
- 🔄 Регулярная ротация ключей доступа
- 📝 Логирование всех API запросов с IP адресами

## 📦 Зависимости

- **Fiber v2** - HTTP веб-фреймворк
- **godotenv** - Загрузка .env файлов
- **AWS SDK** - Интеграция с AWS Secrets Manager
- **MongoDB Driver** - Подключение к MongoDB

🍌 **BELLO!** Простой HTTP API для автоматизации iiko систем! **BANANA!** 🎉