package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	go func() {
		log.Println("gRPC started on", cfg.GRPCPort)
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: ")
		}
	}()

	// Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down....")

	done := make(chan struct{})

	go func() {
		server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Server stopped gracefully")
	case <-time.After(3 * time.Second):
		server.Stop()
	}
}
