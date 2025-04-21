DATABASE_URL := postgres://root:secret@localhost:5432/ranking_psql?sslmode=disable
PATH_MOCK_REPOSITORY := github.com/minhhoanq/video-realtime-ranking/ranking-service/internal/dataaccess

build:
	go build -v ./...

test:
	go test -v ./...

server:
	go run cmd/app/main.go

postgres:
	docker run --name video-realtime-ranking-ranking_psql-1  -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine

createdb:
	docker exec -it video-realtime-ranking-ranking_psql-1 createdb --username root --owner root ranking_psql

dropdb:
	docker exec -it video-realtime-ranking-ranking_psql-1  dropdb ranking_psql

migrate-create:
	migrate create -ext sql -dir ./internal/dataaccess/database/migrations/postgres -seq $(name)

migrateup:
	migrate -path ./internal/dataaccess/database/migrations/postgres -database ${DATABASE_URL} up

migratedown:
	migrate -path ./internal/dataaccess/database/migrations/postgres -database ${DATABASE_URL} down

mockrd:
	mockgen -package mockrd -destination internal/dataaccess/redis/mock/store.go ${PATH_MOCK_REPOSITORY}/redis RankingDataAccessor

.PHONY: build test server up down postgres createdb dropdb migrate-create migrateup migratedown mockrd