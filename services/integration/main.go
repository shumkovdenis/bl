package main

import (
	"log"

	"github.com/shumkovdenis/services/integration/helpers"
)

/*
func newError(transactionID string) error {
	err := connect.NewError(
		connect.CodeInvalidArgument,
		errors.New("player id is required"),
	)

	rollbackInfo := &integration.RollbackInfo{
		TransactionId: transactionID,
	}

	if detail, detailErr := connect.NewErrorDetail(rollbackInfo); detailErr == nil {
		err.AddDetail(detail)
	}

	return err
}
*/

type DaprConfig struct {
	HTTPPort int `env:"HTTP_PORT" envDefault:"3500"`
	GRPCPort int `env:"GRPC_PORT" envDefault:"50001"`
}

type Config struct {
	Dapr      DaprConfig          `envPrefix:"DAPR_"`
	Port      int                 `env:"PORT" envDefault:"6000"`
	Mode      string              `env:"MODE" envDefault:"http"`
	HTTPTrace helpers.TraceConfig `envPrefix:"HTTP_TRACE_"`
	GRPCTrace helpers.TraceConfig `envPrefix:"GRPC_TRACE_"`
}

func main() {
	var cfg Config
	if err := helpers.ParseConfig(&cfg); err != nil {
		log.Fatal(err)
	}

	log.Printf("server started on port %d in %s mode", cfg.Port, cfg.Mode)

	if cfg.Mode == "grpc" {
		if err := NewGRPCServer(cfg); err != nil {
			log.Fatal(err)
		}
	} else if cfg.Mode == "connect" {
		if err := NewConnectServer(cfg); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := NewHTTPServer(cfg); err != nil {
			log.Fatal(err)
		}
	}
}
