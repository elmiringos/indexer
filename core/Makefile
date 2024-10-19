run:
	@go run cmd/api/main.go
.PHONY: run

test:
	@echo "Testing..."
	@go test ./... -v
.PHONY: test

build: 
	@echo "Building..."
	@go build -o main cmd/api/main.go
.PHONY: build

migrate-create:
	@migrate create -ext sql -dir migrations '$(MIGRATE_NAME)'
.PHONY: migrate-create

migrate-up:
	@migrate -path migrations -database '$(PG_URL)' up
.PHONY: migrate-up

migrate-up-c:
	@migrate -path migrations -database '$(PG_URL)' up '$(COUNT)'
.PHONY: migrate-up-c

migrate-down:
	@migrate -path migrations -database '$(PG_URL)' down
.PHONY: migrate-down

migrate-down-c:
	@migrate -path migrations -database '$(PG_URL)' down $(COUNT)
.PHONY: migrate-down-c
