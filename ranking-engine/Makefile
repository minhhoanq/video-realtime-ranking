DATABASE_URL := postgres://root:secret@localhost:5432/ranking_psql?sslmode=disable

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
	migrate create -ext cql -dir ./internal/dataaccess/database/migrations/cassandra -seq $(name)

migrateup:
	migrate -path ./internal/dataaccess/database/migrations/postgres -database ${DATABASE_URL} up

migratedown:
	migrate -path ./internal/dataaccess/database/migrations/postgres -database ${DATABASE_URL} down


.PHONY: build test server up down postgres createdb dropdb migrate-create migrateup migratedown