package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort string
}

func Load() *Config {
	_ = godotenv.Load(".env")
	port := os.Getenv("GRPC_PORT")
	return &Config{GRPCPort: port}
}
