# Makefile
postgres-setup:
	docker-compose up -d

postgres-destroy:
	docker-compose stop
	docker-compose rm -f
	rm -rf ~/simple-bank/postgres

createdb:
	docker exec -it simple-bank-db-1  createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it simple-bank-db-1  dropdb simple_bank

migrate-up:
	migrate -path db/migration -database "postgresql://root:admin@localhost:5433/simple_bank?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://root:admin@localhost:5433/simple_bank?sslmode=disable" -verbose down 

sqlc:
	sqlc generate 

test:
	go test -v -cover ./...

console:
	docker exec -it simple-bank-db-1 psql -U root -d simple_bank

.PHONY: createdb dropdb postgres-destroy postgres-setup migrate-up migrate-down sqlc test