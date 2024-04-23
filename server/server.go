package main

import (
	"context"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/go-microservice-grpc/helloworld"
	"google.golang.org/grpc"
)

type server struct {
	helloworld.HelloworldServiceServer
}

func (s *server) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloResponse, error) {
	return &helloworld.HelloResponse{Message: "Hello " + req.Name}, nil
}

func (s *server) SayHelloServerStream(req *helloworld.HelloRequest, stream helloworld.HelloworldService_SayHelloServerStreamServer) error {
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		if err := stream.Send(&helloworld.HelloResponse{Message: "Hello " + req.Name + " " + strconv.Itoa(i)}); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) SayHelloClientStream(stream helloworld.HelloworldService_SayHelloClientStreamServer) error {
	var names []string
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		names = append(names, req.GetName())
	}

	// Process names received from client
	for _, name := range names {
		time.Sleep(time.Second)
		log.Println(name)
	}

	return stream.SendAndClose(&helloworld.HelloResponse{Message: "done"})
}

func (s *server) SayHelloBidirectionalStream(stream helloworld.HelloworldService_SayHelloBidirectionalStreamServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		name := req.GetName()
		log.Println(name)
		if err := stream.Send(&helloworld.HelloResponse{Message: "Hello " + name}); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Server is running on port 50051")
	s := grpc.NewServer()
	helloworld.RegisterHelloworldServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
