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

PROTO_SRC = internal/api/proto/explorer.proto
PROTO_OUT = internal/api/pb

proto:
	protoc \
	--proto_path=proto \
	--go_out=internal/api \
	--go-grpc_out=internal/api \
	proto/explorer.proto

