package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-microservice-grpc/proto/package/binance_service"
	"github.com/go-microservice-grpc/proto/package/hello_service"
	"google.golang.org/grpc"
)

func main() {
	// Connect to gRPC servers
	connBinance, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect to BinanceService: %v", err)
	}
	defer connBinance.Close()
	binanceClient := binance_service.NewBinanceServiceClient(connBinance)

	connHello, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect to HelloService: %v", err)
	}
	defer connHello.Close()
	helloClient := hello_service.NewHelloServiceClient(connHello)

	// Initialize Gin
	r := gin.Default()

	// Define routes
	r.GET("/api/v1/klines/:symbol", func(c *gin.Context) {
		symbol := c.Param("symbol")
		response, err := binanceClient.GetKlinesBySymbol(context.Background(), &binance_service.GetKlinesBySymbolRequest{
			Symbol:   symbol,
			EndTime:  time.Now().UnixMilli(),
			Limit:    1000,
			Interval: "1d",
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, response.Data)
	})

	r.GET("/api/v1/hello/:name", func(c *gin.Context) {
		name := c.Param("name")
		response, err := helloClient.SayHello(context.Background(), &hello_service.HelloRequest{Name: name})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": response.Message})
	})

	r.GET("/api/v1/hello_stream/:name", func(c *gin.Context) {
		name := c.Param("name")
		stream, err := helloClient.SayHelloServerStream(context.Background(), &hello_service.HelloRequest{Name: name})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		messages := []string{}
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			messages = append(messages, response.Message)
		}
		c.JSON(http.StatusOK, gin.H{"messages": messages})
	})

	r.POST("/api/v1/hello_stream", func(c *gin.Context) {
		var names []string
		if err := c.ShouldBindJSON(&names); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		clientStream, err := helloClient.SayHelloClientStream(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, name := range names {
			if err := clientStream.Send(&hello_service.HelloRequest{Name: name}); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		response, err := clientStream.CloseAndRecv()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": response.Message})
	})

	r.GET("/api/v1/hello_bidi", func(c *gin.Context) {
		bidiStream, err := helloClient.SayHelloBidirectionalStream(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var wgBidi sync.WaitGroup
		wgBidi.Add(1)
		messages := make(chan string)

		go func() {
			defer wgBidi.Done()
			for {
				response, err := bidiStream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				messages <- response.Message
			}
			close(messages)
		}()

		for i := 0; i < 5; i++ {
			time.Sleep(time.Second)
			if err := bidiStream.Send(&hello_service.HelloRequest{Name: "David " + strconv.Itoa(i)}); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		bidiStream.CloseSend()
		wgBidi.Wait()

		var result []string
		for msg := range messages {
			result = append(result, msg)
		}

		c.JSON(http.StatusOK, gin.H{"messages": result})
	})

	// Run the Gin server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
