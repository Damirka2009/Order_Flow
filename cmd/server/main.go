package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"master/internal/config"
	"master/internal/repository"
	"master/internal/service"

	grpcHandler "master/internal/transport/grpc"
	api "master/pkg/api"
)

func main() {
	cfg := config.Load()
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.New()
	svc := service.New(repo)
	handler := grpcHandler.New(svc)

	server := grpc.NewServer()

	api.RegisterOrderServiceServer(server, handler)
	log.Println("gRPC started on ", cfg.GRPCPort)
	server.Serve(lis)
}
