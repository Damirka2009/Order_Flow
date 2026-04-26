package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"master/internal/inventory"
	grpcHandler "master/internal/inventory/transport/grpc"
	api "master/pkg/api"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}

	svc := inventory.New()
	handler := grpcHandler.New(svc)

	server := grpc.NewServer()
	api.RegisterInventoryServiceServer(server, handler)

	log.Println("Inventory service running on :50052")

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
