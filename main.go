package main

import (
	"github.com/rs/zerolog/log"
)

func main() {
	var cfg Config
	if err := ParseConfig(&cfg); err != nil {
		log.Fatal().Err(err).Msg("failed to parse config")
	}
	cfg.Log()

	var callee Callee

	if cfg.Callee.Mode == "http" {
		callee = NewHTTPCallee(cfg)
	} else if cfg.Callee.Mode == "connect" {
		callee = NewConnectCallee(cfg)
	} else if cfg.Callee.Mode == "error" {
		callee = NewErrorCallee()
	} else if cfg.Callee.Mode == "grpc" {
		callee = NewGRPCCallee(cfg)
	} else {
		callee = NewFakeCallee()
	}

	if cfg.Mode == "http" {
		if err := NewHTTPService(cfg, callee); err != nil {
			log.Fatal().Err(err).Msg("failed to start http server")
		}
	} else if cfg.Mode == "connect" {
		if err := NewConnectService(cfg, callee); err != nil {
			log.Fatal().Err(err).Msg("failed to start connect server")
		}
	} else if cfg.Mode == "grpc" {
		if err := NewGRPCService(cfg, callee); err != nil {
			log.Fatal().Err(err).Msg("failed to start grpc server")
		}
	} else {
		log.Fatal().Msg("unknown mode")
	}
}
