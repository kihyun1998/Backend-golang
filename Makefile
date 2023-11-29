DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres15 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres15 dropdb simple_bank

migrateup:
	migrate -path db/migration/ -database "$(DB_URL)" -verbose up

migrateup_onestep:
	migrate -path db/migration/ -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration/ -database "$(DB_URL)" -verbose down

migratedown_onestep:
	migrate -path db/migration/ -database "$(DB_URL)" -verbose down 1


sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go simplebank/db/sqlc Store

db_docs_build:
	dbdocs build doc/db.dbml

create_db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

proto:
	del /s .\pb\*.go &
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/*.proto

.PHONY: postgres createdb dropdb migrateup migratedown migrateup_onestep migratedown_onestep sqlc test server mock db_docs_build create_db_schema proto