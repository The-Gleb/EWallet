postgres:
	docker run --name emoney_db -e POSTGRES_USER=emoney -e POSTGRES_PASSWORD=emoney -p 5434:5432 -d postgres:alpine
postgresrm:
	docker stop emoney_db
	docker rm emoney_db

migrateup:
	migrate -path internal/adapter/db/postgres/migration -database "postgres://emoney:emoney@localhost:5434/emoney?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/adapter/db/postgres/migration -database "postgres://emoney:emoney@localhost:5434/emoney?sslmode=disable" -verbose down
.PHONY: postgres postgresrm migrateup migratedown