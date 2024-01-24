.PHONY: swag-install
swag-install:
	@go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: swag
swag:
	@`go env GOPATH`/bin/swag init --parseDependency --parseInternal

.PHONY: build
build:
	@go build -o ./build/app ./main.go

.PHONY: run
run: swag build
	@./build/app