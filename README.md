# PR-Reviewer-Assignment-Service

Микросервис для автоматического назначения ревьюеров на Pull Request'ы с управлением командами и пользователями. 

## Описание проекта

Сервис предоставляет REST API для:
- Управления командами разработчиков и их участниками
- Автоматического назначения ревьюеров на Pull Request'ы
- Переназначения ревьюеров
- Мержа Pull Request'ов
- Получения статистики назначений

## Быстрый старт

### Через Docker Compose (рекомендуется)

```bash
git clone <repository-url>
cd pr-reviewer-assignment-service

make docker-compose-up

# Сервис будет доступен на http://localhost:8080
```

### Локальный запуск

```bash
make deps

docker run -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:15-alpine

# Создать базу данных
createdb -U postgres -h localhost reviewer_assigner

# Запустить миграции
make migrate-up

# Запустить сервис
make run
```

## API Документация

Полная спецификация API доступна в файле [`openapi.yml`](./openapi.yml).

### Основные эндпоинты

#### Команды
- `POST /api/team/add` - Создание команды с участниками
- `GET /api/team/get?team_name={name}` - Получение информации о команде

#### Пользователи
- `POST /api/users/setIsActive` - Изменение статуса активности пользователя (требует admin токена)
- `GET /api/users/getReview?user_id={id}` - Получение PR для ревьювера

#### Pull Requests
- `POST /api/pullRequest/create` - Создание PR с автоматическим назначением ревьюверов
- `POST /api/pullRequest/merge` - Мерж PR
- `POST /api/pullRequest/reassign` - Переназначение ревьювера

#### Проверка состояния
- `GET /health` - Проверка здоровья сервиса

### Аутентификация

API использует токены для аутентификации:
- `admin-token` - для операций требующих прав администратора
- `user-token` - для обычных операций

Токены передаются в заголовке `Authorization` или параметре `token`.

## Архитектурные решения

### Общая архитектура

```
├── cmd/server - Точка входа приложения
├── internal/
│   ├── config - Конфигурация приложения
│   ├── database - Подключение к БД
│   ├── models - Структуры данных
│   ├── handlers - HTTP обработчики
│   ├── services - Бизнес-логика
│   ├── repository - Работа с БД
│   └── middleware - Middleware компоненты
├── migrations - SQL миграции
└── tests - Тесты (unit, integration, e2e)
```

### Ключевые решения

#### 1. Идемпотентность операций

- Операция merge реализована идемпотентно - повторный вызов не вызывает ошибки
- Все операции валидируют состояние перед выполнением

#### 2. Безопасность данных

- Использование транзакций для операций изменения ревьюверов
- Строгая валидация состояния PR перед модификацией
- Каскадные обновления при изменении команд

#### 3. Масштабируемость

- Репозиторий паттерн для абстракции работы с БД
- Сервисный слой для бизнес-логики
- Поддержка health checks для оркестрации

### Используемые технологии

- **Go 1.25.1** - основной язык программирования
- **Gin** - HTTP фреймворк
- **PostgreSQL 15** - база данных
- **pgx/v5** - PostgreSQL драйвер
- **golang-migrate** - управление миграциями
- **testify** - фреймворк для тестирования
- **testcontainers** - интеграционное тестирование
- **Docker & Docker Compose** - контейнеризация

### Конфигурация

Приложение настраивается через переменные окружения:

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| `SERVER_PORT` | Порт сервера | 8080 |
| `SERVER_READ_TIMEOUT` | Таймаут чтения (сек) | 30 |
| `SERVER_WRITE_TIMEOUT` | Таймаут записи (сек) | 30 |
| `DB_HOST` | Хост БД | localhost |
| `DB_PORT` | Порт БД | 5432 |
| `DB_USER` | Пользователь БД | your-user |
| `DB_PASSWORD` | Пароль БД | your-password |
| `DB_NAME` | Имя БД | your-db |
| `DB_SSLMODE` | Режим SSL | disable |

### Запуск тестов

```bash
# Все тесты
make test

# С покрытием
make test-coverage

# Только unit тесты
make test-unit

# Интеграционные тесты
make test-integration

# E2E тесты
make test-e2e
```

## Разработка

### Сборка и запуск

```bash
# Сборка бинарного файла
make build

# Запуск
make run

# Очистка
make clean
```

### Работа с миграциями

```bash
# Применить миграции
make migrate-up

# Откатить миграции
make migrate-down

# Создать новую миграцию
make migrate-create name=your_migration_name
```

### Docker команды

```bash
# Сборка образа
make docker-build

# Запуск контейнера
make docker-run

# Docker Compose
make docker-compose-up
make docker-compose-down
make docker-compose-logs
```
