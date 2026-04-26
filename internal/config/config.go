package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found (ok in prod)")
	}
	port := os.Getenv("GRPC_PORT")

	return &Config{GRPCPort: port}
}
