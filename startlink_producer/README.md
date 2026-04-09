# starlink_producer

## Запуск

### 1. Настройка окружения

Создайте файл `.env` в корне проекта и вставить туда содержимое .env.example

### 2. Миграции
Установить мигратор:
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Провести миграцию:
```bash
make migrate-up
```

### 2. Установка зависимостей

```bash
go mod tidy
```

### 3. Запуск

```bash
go run ./cmd/api/
```

### 4. Сборка и запуск бинарника

```bash
go build -o starlink_producer ./cmd/api/
./starlink_producer
```

## API

### Создать пользователя

```
POST /api/users
Content-Type: application/json

{
  "first_name": "Ivan",
  "last_name": "Ivanov",
  "email": "ivan@example.com"
}
```
