postgres:
	sudo docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:17.5-alpine3.22

createdb: 
	sudo docker exec -it postgres17 createdb --username=root --owner=root gobank

dropdb: 
	sudo docker exec -it postgres17 dropdb gobank

migrateup:
	./migrate -path db/migration -database "postgresql://root:123456@localhost:5432/gobank?sslmode=disable" -verbose up

migratedown:
	./migrate -path db/migration -database "postgresql://root:123456@localhost:5432/gobank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc