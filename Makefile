gen:
	protoc --go_out=. --go-grpc_out=. proto/*.proto
gen-pb:
	protoc --go_out=. proto/*.proto
gen-grpc-pb:
	protoc --go-grpc_out=. proto/*.proto

ser:
	go run ./server/server.go
cli:
	go run ./client/client.go