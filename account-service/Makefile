.PHONY: swag-install
swag-install:
	@go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: swag
swag:
	@`go env GOPATH`/bin/swag init --parseDependency --parseInternal

.PHONY: proto
proto:
	@protoc -I=. --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. grpc/user.proto

.PHONY: build
build:
	@go build -o ./build/app ./main.go

.PHONY: run
run: swag build
	@./build/app