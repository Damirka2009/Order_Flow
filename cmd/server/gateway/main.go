package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"master/internal/config"
	api "master/pkg/api"
)

func main() {
	cfg := config.Load()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	port := "localhost:" + cfg.GRPCPort

	err := api.RegisterOrderServiceHandlerFromEndpoint(
		ctx,
		mux,
		port,
		opts,
	)
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	go func() {
		log.Println("HTTP Gateway started on :8080")
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("failed to start gateway: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("Studding down gateway....")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatal("Gateway forced to shutdown:", err)
	}
	log.Println("Gateway stopped")
}
