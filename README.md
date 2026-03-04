```markdown
# Fitness Club Energy

Веб-приложение для фитнес-клуба на Go с PostgreSQL и Redis.

## Функциональность

- Главная страница с информацией о клубе
- Страница регистрации клиентов
- Страница входа в систему
- Просмотр абонементов
- Информация о клубе и контакты
- Управление клиентами (CRUD операции)
- Система абонементов
- Учет тренировок
- Управление тренерами

## Технологии

- **Backend:** Go 1.21+
- **Базы данных:** PostgreSQL, Redis (кэширование)
- **Фреймворки:** Gin (HTTP), GORM (ORM)
- **Логирование:** Uber Zap
- **Инфраструктура:** Docker, Docker Compose

## Запуск проекта

### 1. Запуск базы данных и Redis
```bash
docker-compose up -d
```

### 2. Запуск приложения
```bash
go mod tidy
go run cmd/main.go
```

### 3. Доступ к приложению
- **Сервер:** `http://localhost:8080`
- **Swagger документация:** `http://localhost:8080/swagger/index.html`

## Структура проекта

```
.
├── cmd/
│   └── main.go                 # Точка входа
├── internal/
│   ├── cache/                   # Redis клиент и обёртки
│   │   ├── redis.go
│   │   ├── wrapper.go
│   │   └── keys.go
│   ├── config/                  # Конфигурация
│   │   └── config.go
│   ├── dto/                     # Data Transfer Objects
│   │   ├── request/             # Входящие DTO
│   │   └── response/            # Исходящие DTO
│   ├── handler/                  # HTTP обработчики
│   │   ├── user.go
│   │   ├── membership.go
│   │   ├── workout.go
│   │   └── router.go
│   ├── logger/                   # Логирование (zap)
│   │   └── logger.go
│   ├── model/                    # Модели данных
│   │   ├── user.go
│   │   ├── membership.go
│   │   └── workout.go
│   ├── repository/                # Работа с БД
│   │   ├── database.go
│   │   ├── user_repository.go
│   │   ├── membership_repository.go
│   │   └── workout_repository.go
│   └── service/                   # Бизнес-логика
│       ├── user_service.go
│       ├── membership_service.go
│       └── workout_service.go
├── migrations/                    # SQL миграции
│   └── 001_init.sql
├── docker-compose.yml
├── .env.example
├── go.mod
└── README.md
```

## Переменные окружения

Создайте файл `.env` на основе `.env.example`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres123
DB_NAME=fitness_club
SERVER_PORT=8080
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

## Автор

beushvvv
```