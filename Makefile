gen:
	protoc --go_out=. --go-grpc_out=. proto/*.proto

build-binance:
	@go build -o bin/binance_service ./server/binance_service

build-hello:
	@go build -o bin/hello_service ./server/hello_service

build:
	@go build -o bin/main

run-binance: build-binance
	bin/binance_service

run-hello: build-hello
	bin/hello_service

run: build
	bin/main
