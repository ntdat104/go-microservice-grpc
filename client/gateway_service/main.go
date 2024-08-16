package main

import (
	"context"
	"log"
	"time"

	"github.com/go-microservice-grpc/proto/package/binance_service"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := binance_service.NewBinanceServiceClient(conn)

	response, err := c.GetKlinesBySymbol(context.Background(), &binance_service.GetKlinesBySymbolRequest{
		Symbol:   "BTCUSDT",
		EndTime:  time.Now().UnixMilli(),
		Limit:    1000,
		Interval: "1d",
	})
	if err != nil {
		log.Fatalf("GetKlinesBySymbol has err: %v", err)
	}

	for _, v := range response.Data {
		volume := 1.4
		v.Volume = &volume;
		log.Println(*v.Volume)
	}

	// conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	// defer conn.Close()
	// c := hello_service.NewHelloServiceClient(conn)

	// // Unary RPC
	// response, err := c.SayHello(context.Background(), &hello_service.HelloRequest{Name: "Alice"})
	// if err != nil {
	// 	log.Fatalf("could not greet: %v", err)
	// }
	// log.Printf("Greeting: %s", response.Message)

	// // Server streaming RPC
	// stream, err := c.SayHelloServerStream(context.Background(), &hello_service.HelloRequest{Name: "Bob"})
	// if err != nil {
	// 	log.Fatalf("error while calling SayHelloServerStream: %v", err)
	// }
	// for {
	// 	response, err := stream.Recv()
	// 	if err != nil {
	// 		break
	// 	}
	// 	log.Printf("Greeting from server stream: %s", response.Message)
	// }

	// // Client streaming RPC
	// clientStream, err := c.SayHelloClientStream(context.Background())
	// if err != nil {
	// 	log.Fatalf("error while calling SayHelloClientStream: %v", err)
	// }

	// var names = []string{"Charlie", "Daniel", "Emma", "Fiona", "George"}
	// for _, name := range names {
	// 	if err := clientStream.Send(&hello_service.HelloRequest{Name: name}); err != nil {
	// 		log.Fatalf("error sending request: %v", err)
	// 	}
	// }

	// response, err = clientStream.CloseAndRecv()
	// if err != nil {
	// 	log.Fatalf("error receiving response: %v", err)
	// }
	// log.Printf("Greeting from client stream: %s", response.Message)

	// // Bidirectional streaming RPC
	// bidiStream, err := c.SayHelloBidirectionalStream(context.Background())
	// if err != nil {
	// 	log.Fatalf("error while calling SayHelloBidirectionalStream: %v", err)
	// }
	// var wgBidi sync.WaitGroup
	// wgBidi.Add(1)
	// go func() {
	// 	defer wgBidi.Done()
	// 	for {
	// 		response, err := bidiStream.Recv()
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		if err != nil {
	// 			log.Fatalf("error receiving response: %v", err)
	// 		}
	// 		log.Printf("Greeting from bidirectional stream: %s", response.Message)
	// 	}
	// }()
	// for i := 0; i < 5; i++ {
	// 	time.Sleep(time.Second)
	// 	if err := bidiStream.Send(&hello_service.HelloRequest{Name: "David " + strconv.Itoa(i)}); err != nil {
	// 		log.Fatalf("error sending request: %v", err)
	// 	}
	// }
	// bidiStream.CloseSend()
	// wgBidi.Wait()
}
