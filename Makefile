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

sqlc:
	sqlc generate

test:
	go test -v -cover ./... 


.PHONY: createdb dropdb postgres migrateup migratedown sqlc test
