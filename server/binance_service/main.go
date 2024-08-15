package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/go-microservice-grpc/proto/package/binance_service"
	"google.golang.org/grpc"
)

type server struct {
	binance_service.BinanceServiceServer
}

func (s *server) GetKlinesBySymbol(ctx context.Context, req *binance_service.GetKlinesBySymbolRequest) (*binance_service.GetKlinesBySymbolResponse, error) {
	log.Println("GetKlinesBySymbol:", req)

	url := fmt.Sprintf("https://www.binance.com/api/v3/uiKlines?endTime=%d&limit=%d&symbol=%s&interval=%s", req.EndTime, req.Limit, req.Symbol, req.Interval)

	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to make request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return nil, err
	}

	// Parse the JSON response
	var klines [][]interface{}
	if err := json.Unmarshal(body, &klines); err != nil {
		fmt.Printf("Failed to parse JSON: %v\n", err)
		return nil, err
	}

	var parsedKlines []*binance_service.KlineData
	for _, k := range klines {
		open, err := strconv.ParseFloat(k[1].(string), 64)
		if err != nil {
			log.Fatalf("Error converting Open: %v", err)
		}

		high, err := strconv.ParseFloat(k[2].(string), 64)
		if err != nil {
			log.Fatalf("Error converting High: %v", err)
		}

		low, err := strconv.ParseFloat(k[3].(string), 64)
		if err != nil {
			log.Fatalf("Error converting Low: %v", err)
		}

		closePrice, err := strconv.ParseFloat(k[4].(string), 64)
		if err != nil {
			log.Fatalf("Error converting Close: %v", err)
		}

		volume, err := strconv.ParseFloat(k[5].(string), 64)
		if err != nil {
			log.Fatalf("Error converting Volume: %v", err)
		}
		kline := &binance_service.KlineData{
			OpenTime:  int64(k[0].(float64)),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     closePrice,
			Volume:    volume,
			CloseTime: int64(k[6].(float64)),
		}
		parsedKlines = append(parsedKlines, kline)
	}

	return &binance_service.GetKlinesBySymbolResponse{Data: parsedKlines}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Binance server is running on port 50052")
	s := grpc.NewServer()
	binance_service.RegisterBinanceServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
