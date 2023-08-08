createdb:
	docker exec -it postgres23 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -t postgres23 dropdb simple_bank


postgres:
	docker run --name postgres23 -p 5432:5432  -e POSTGRES_USER=root -e POSTGRES_PASSWORD=admin -d postgres:15-alpine

migrateup:
	migrate -path ./db/migration/ -database "postgresql://root:admin@34.206.16.110:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path ./db/migration/ -database "postgresql://root:admin@34.206.16.110:5432/simple_bank?sslmode=disable" -verbose down -all

migrateup1:
	migrate -path ./db/migration/ -database "postgresql://root:admin@34.206.16.110:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path ./db/migration/ -database "postgresql://root:admin@34.206.16.110:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./... 

server:
	go run main.go


.PHONY: createdb dropdb postgres migrateup migratedown migrateup1 migratedown1 sqlc test server
