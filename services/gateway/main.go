package main

import (
	"log"

	"github.com/shumkovdenis/bl/services/gateway/helpers"
)

/*
func extractError(err error) (*integration.RollbackInfo, bool) {
	var connectErr *connect.Error
	if !errors.As(err, &connectErr) {
		return nil, false
	}

	for _, detail := range connectErr.Details() {
		msg, valueErr := detail.Value()
		if valueErr != nil {
			// Usually, errors here mean that we don't have the schema for this
			// Protobuf message.
			continue
		}
		if retryInfo, ok := msg.(*integration.RollbackInfo); ok {
			return retryInfo, true
		}
	}

	return nil, false
}
*/

type DaprConfig struct {
	HTTPPort int `env:"HTTP_PORT" envDefault:"3500"`
	GRPCPort int `env:"GRPC_PORT" envDefault:"50001"`
}

type IntegrationConfig struct {
	AppID string `env:"APP_ID" envDefault:"integration"`
}

type Config struct {
	Dapr        DaprConfig          `envPrefix:"DAPR_"`
	Port        int                 `env:"PORT" envDefault:"6000"`
	HTTPTrace   helpers.TraceConfig `envPrefix:"HTTP_TRACE_"`
	GRPCTrace   helpers.TraceConfig `envPrefix:"GRPC_TRACE_"`
	Integration IntegrationConfig   `envPrefix:"INTEGRATION_"`
}

func main() {
	var cfg Config
	if err := helpers.ParseConfig(&cfg); err != nil {
		log.Fatal(err)
	}

	log.Printf("server started on port %d", cfg.Port)

	if err := NewHTTPServer(cfg); err != nil {
		log.Fatal(err)
	}
}
