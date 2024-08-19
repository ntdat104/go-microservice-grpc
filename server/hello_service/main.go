package main

import (
	"context"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	pb "github.com/go-microservice-grpc/proto/package/hello_service"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedHelloServiceServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello " + req.Name}, nil
}

func (s *server) SayHelloServerStream(req *pb.HelloRequest, stream pb.HelloService_SayHelloServerStreamServer) error {
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		if err := stream.Send(&pb.HelloResponse{Message: "Hello " + req.Name + " " + strconv.Itoa(i)}); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) SayHelloClientStream(stream pb.HelloService_SayHelloClientStreamServer) error {
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

	return stream.SendAndClose(&pb.HelloResponse{Message: "done"})
}

func (s *server) SayHelloBidirectionalStream(stream pb.HelloService_SayHelloBidirectionalStreamServer) error {
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
		if err := stream.Send(&pb.HelloResponse{Message: "Hello " + name}); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Server is running on port: 50051")
	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
