package main

import (
	"log"
)

func main() {
	var cfg Config
	if err := ParseConfig(&cfg); err != nil {
		log.Fatal("failed to parse config:", err)
	}
	cfg.Log()

	var callee Callee

	if cfg.Callee.Mode == "http" {
		callee = NewHTTPCallee(cfg)
	} else if cfg.Callee.Mode == "connect" {
		callee = NewConnectCallee(cfg)
	} else {
		callee = NewFakeCallee()
	}

	if cfg.Mode == "http" {
		if err := NewHTTPService(cfg, callee); err != nil {
			log.Fatal("failed to start http server:", err)
		}
	} else if cfg.Mode == "connect" {
		if err := NewConnectService(cfg, callee); err != nil {
			log.Fatal("failed to start connect server:", err)
		}
	} else {
		log.Fatal("unknown mode")
	}
}
