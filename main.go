package main

import "log"

func main() {
	var cfg Config
	if err := ParseConfig(&cfg); err != nil {
		log.Fatal(err)
	}
	cfg.Log()

	if err := NewHTTPServer(cfg, NewFakeCaller()); err != nil {
		log.Fatal(err)
	}
}
