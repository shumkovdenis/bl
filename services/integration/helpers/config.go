package helpers

import "github.com/caarlos0/env/v8"

type TraceConfig struct {
	UseTraceParentHeader  bool `env:"TRACE_PARENT" envDefault:"false"`
	UseTraceStateHeader   bool `env:"TRACE_STATE" envDefault:"false"`
	UseGrpcTraceBinHeader bool `env:"GRPC_TRACE_BIN" envDefault:"false"`
}

func ParseConfig(cfg interface{}) error {
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return err
	}
	return nil
}
