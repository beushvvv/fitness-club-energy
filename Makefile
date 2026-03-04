.PHONY: run build swag docker-up docker-down migrate
.PHONY: git-branch git-commit git-push

git-branch:
	git checkout -b feature/swagger-dto

git-commit:
	git add .
	git commit -m "feat: add DTO models, Swagger documentation, dynamic data handling"

git-push:
	git push origin feature/swagger-dto

run:
	go run cmd/main.go

build:
	go build -o fitness-club cmd/main.go

swag:
	swag init -g cmd/main.go -o docs

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate:
	docker exec -it fitness-postgres psql -U postgres -d fitness_club -f /docker-entrypoint-initdb.d/001_create_users.sql

test:
	curl http://localhost:8080/api/v1/users
	curl http://localhost:8080/api/v1/memberships