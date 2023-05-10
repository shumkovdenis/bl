package helpers

import "github.com/caarlos0/env/v8"

func ParseConfig(cfg interface{}) error {
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return err
	}
	return nil
}
