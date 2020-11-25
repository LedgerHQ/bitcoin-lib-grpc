package main

import (
	"fmt"
	"net"

	"github.com/ledgerhq/bitcoin-lib-grpc/config"

	controllers "github.com/ledgerhq/bitcoin-lib-grpc/grpc"
	"github.com/ledgerhq/bitcoin-lib-grpc/log"
	pb "github.com/ledgerhq/bitcoin-lib-grpc/pb/bitcoin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func serve(addr string) {
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Cannot listen to address %s", addr)
	}

	s := grpc.NewServer()
	bitcoinController := controllers.NewBitcoinController()
	healthController := controllers.NewHealthChecker()

	pb.RegisterCoinServiceServer(s, bitcoinController)
	grpc_health_v1.RegisterHealthServer(s, healthController)

	reflection.Register(s)

	if err := s.Serve(conn); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func main() {
	configProvider := config.LoadProvider("bitcoin")

	var (
		host string
		port int32 = 50051
	)

	host = configProvider.GetString("host")

	if val := configProvider.GetInt32("port"); val != 0 {
		port = val
	}

	addr := fmt.Sprintf("%s:%d", host, port)

	serve(addr)
}
