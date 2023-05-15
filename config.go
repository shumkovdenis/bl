package main

import (
	"log"

	"github.com/caarlos0/env/v8"
)

type DaprConfig struct {
	HTTPPort int `env:"HTTP_PORT" envDefault:"3500"`
	GRPCPort int `env:"GRPC_PORT" envDefault:"50001"`
}

type CalleeConfig struct {
	ServiceName string `env:"SERVICE_NAME" envDefault:"service"`
	Mode        string `env:"MODE" envDefault:"http"`
}

type Config struct {
	Dapr        DaprConfig   `envPrefix:"DAPR_"`
	ServiceName string       `env:"SERVICE_NAME" envDefault:"service"`
	Mode        string       `env:"MODE" envDefault:"http"`
	Port        int          `env:"PORT" envDefault:"6000"`
	Callee      CalleeConfig `envPrefix:"CALLEE_"`
}

func (c *Config) Log() {
	log.Println("service name:", c.ServiceName)
	log.Println("mode:", c.Mode)
	log.Println("port:", c.Port)
	log.Println("callee service name:", c.Callee.ServiceName)
	log.Println("callee mode:", c.Callee.Mode)
}

func ParseConfig(cfg interface{}) error {
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return err
	}
	return nil
}
