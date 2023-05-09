package main

import (
	"log"
)

// type IntegrationConfig struct {
// 	AppID string `env:"APP_ID" envDefault:"integration"`
// }

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

func main() {
	cfg, err := ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := NewHTTPServer(cfg); err != nil {
		log.Fatal(err)
	}

	log.Printf("server started on port %d in", cfg.Port)
}
