package config

import "github.com/kelseyhightower/envconfig"

func NewConfig() (cfg Config, err error) {
	err = envconfig.Process("", &cfg)
	cfg.ServiceName = "org"
	return
}
