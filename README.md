# 🍌 Minion - Автоматизация iiko API

BELLO! Минион для автоматизации работы с iiko API системами. Продление ключей и обновление меню для ресторанов.

## ✨ Возможности

- 🔑 **Продление API ключей** - автоматическое продление срока действия API ключей до максимального значения
- 🍽️ **Обновление меню** - обновление внешних меню для всех ресторанов
- 🔐 **AWS Secrets Manager** - безопасное хранение данных подключения к базе данных
- 📊 **MongoDB интеграция** - загрузка ресторанов из базы данных
- 🎯 **Фильтрация** - обработка только активных iiko ресторанов
- 🍌 **Минион тематика** - BELLO! BANANA! GELATO! POOPAYE!

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
go mod download

# Создаем .env файл
cp env.example .env

# Редактируем .env файл
nano .env
```

### Настройка

Создайте `.env` файл на основе `env.example`:

```bash
# AWS Secrets Manager
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
make build

# Продление API ключей
./bin/minion extend-keys

# Обновление меню
./bin/minion refresh-menus
# Помощь
./bin/minion --help
```

## 🔧 Конфигурация

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
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

## 🏗️ Архитектура

```
cmd/minion/           - Точка входа приложения
internal/
├── aws/             - AWS Secrets Manager клиент
├── client/          - HTTP клиент для iiko API
├── commands/        - Cobra команды (extend-keys, refresh-menus)
├── config/          - Конфигурация и загрузка ресторанов
├── database/        - MongoDB сервис
└── models/          - Структуры данных
```

## 🛠️ Команды Make

```bash
make build          # Сборка бинарного файла
make deps           # Установка зависимостей
make help           # Показать все команды
```

## 📦 Зависимости

- **Cobra** - CLI фреймворк
- **godotenv** - Загрузка .env файлов
- **AWS SDK** - Интеграция с AWS Secrets Manager
- **MongoDB Driver** - Подключение к MongoDB

---

🍌 **BELLO!** Сделано с любовью для автоматизации iiko систем! **BANANA!** 🎉