package main

import (
	"github.com/caarlos0/env/v8"
)

type DaprConfig struct {
	HHTPPort int `env:"HTTP_PORT" envDefault:"3500"`
	GRPCPort int `env:"GRPC_PORT" envDefault:"50001"`
}

type Config struct {
	Dapr DaprConfig `envPrefix:"DAPR_"`
	Port int        `env:"PORT" envDefault:"6000"`
}

func ParseConfig() (Config, error) {
	var cfg Config
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
