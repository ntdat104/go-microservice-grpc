gen-core-service:
	protoc --go_out=. --go-grpc_out=. proto/core_service/*.proto

gen-hello-service:
	protoc --go_out=. --go-grpc_out=. proto/hello_service/*.proto