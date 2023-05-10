package main

import (
	"log"

	"github.com/shumkovdenis/services/remote/helpers"
)

type DaprConfig struct {
	HTTPPort int `env:"HTTP_PORT" envDefault:"3500"`
	GRPCPort int `env:"GRPC_PORT" envDefault:"50001"`
}

type Config struct {
	Dapr DaprConfig `envPrefix:"DAPR_"`
	Port int        `env:"PORT" envDefault:"6000"`
}

func main() {
	var cfg Config
	if err := helpers.ParseConfig(&cfg); err != nil {
		log.Fatal(err)
	}

	log.Printf("server started on port %d", cfg.Port)

	if err := NewGRPCServer(cfg); err != nil {
		log.Fatal(err)
	}
}
