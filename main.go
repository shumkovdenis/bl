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
	} else {
		callee = NewFakeCallee()
	}

	if err := NewHTTPService(cfg, callee); err != nil {
		log.Fatal("failed to start http server:", err)
	}
}
