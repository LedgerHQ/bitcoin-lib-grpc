package main

import (
	"fmt"
	"net"

	"github.com/ledgerhq/lama-bitcoin-svc/log"
	pb "github.com/ledgerhq/lama-bitcoin-svc/pb/v1"
	"github.com/ledgerhq/lama-bitcoin-svc/pkg/bitcoin"
	"google.golang.org/grpc"
)

func serve() {
	addr := fmt.Sprintf(":%d", 50051)

	conn, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Cannot listen to address %s", addr)
	}

	s := grpc.NewServer()
	pb.RegisterCoinServiceServer(s, &bitcoin.Service{})

	if err := s.Serve(conn); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func main() {
	serve()
}
