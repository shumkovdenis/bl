package main

import (
	"log"

	"github.com/caarlos0/env/v8"
)

type DaprConfig struct {
	HTTPPort int `env:"HTTP_PORT" envDefault:"3500"`
	GRPCPort int `env:"GRPC_PORT" envDefault:"50001"`
}

type CallServiceConfig struct {
	Name string `env:"NAME"`
}

type Config struct {
	Dapr        DaprConfig        `envPrefix:"DAPR_"`
	ServiceName string            `env:"SERVICE_NAME"`
	Mode        string            `env:"MODE" envDefault:"http"`
	Port        int               `env:"PORT" envDefault:"6000"`
	CallService CallServiceConfig `envPrefix:"CALL_SERVICE_"`
}

func (c *Config) Log() {
	log.Println("service name:", c.ServiceName)
	log.Println("mode:", c.Mode)
	log.Println("port:", c.Port)
	log.Println("call service name:", c.CallService.Name)
}

func ParseConfig(cfg interface{}) error {
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return err
	}
	return nil
}
