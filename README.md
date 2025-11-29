# Fitness Club Management System

Система управления фитнес клубом на Go с PostgreSQL

## Функциональность

- Управление клиентами (CRUD операции)
- Система абонементов
- Учет тренировок
- Управление тренерами

## Технологии

- Go 1.21
- PostgreSQL
- GORM ORM
- Gin Web Framework
- Docker

## Запуск

```bash
# Запуск базы данных
docker-compose up -d

# Запуск приложения
go run main.go
