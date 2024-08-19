migrate:
	@migrate create -ext sql -dir db/migration -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@migrate -path db/migration -database "postgres://jesk:testing@localhost:5433/simple_bank_jwt?sslmode=disable&search_path=tutorial" --verbose up 1

migrate-down:
	@migrate -path db/migration -database "postgres://jesk:testing@localhost:5433/simple_bank_jwt?sslmode=disable&search_path=tutorial" --verbose down 1

migrate-all:
	@migrate -path db/migration -database "postgres://jesk:testing@localhost:5433/simple_bank_jwt?sslmode=disable&search_path=tutorial" --verbose up

test:
	go test ./.. -cover -v

sqlc:
	@sqlc generate

.PHONY: migrate migrate-up migrate-down sqlc migrate-all test