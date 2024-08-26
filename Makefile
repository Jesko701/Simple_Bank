migrate:
	@migrate create -ext sql -dir db/migration -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@migrate -path db/migration -database "postgres://jesk:testing@localhost:5433/simple_bank_jwt?sslmode=disable&search_path=tutorial" --verbose up 1

migrate-down:
	@migrate -path db/migration -database "postgres://jesk:testing@localhost:5433/simple_bank_jwt?sslmode=disable&search_path=tutorial" --verbose down 1

migrate-test:
	@migrate -path db/migration -database "$(DATABASE_URL)" --verbose up

test:
	go test ./... -cover -v

server:
	go run main.go

sqlc:
	@sqlc generate

.PHONY: migrate migrate-up migrate-down sqlc migrate-test test server