package main

import (
	"fmt"
	"net"

	protobuf "github.com/ledgerhq/lama-bitcoin-svc"
	"github.com/ledgerhq/lama-bitcoin-svc/log"
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
	protobuf.RegisterBitcoinServiceServer(s, &bitcoin.Service{})

	if err := s.Serve(conn); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func main() {
	serve()
}
