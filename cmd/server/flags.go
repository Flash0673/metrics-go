package main

import (
	"github.com/caarlos0/env/v11"
	"github.com/spf13/pflag"
)

var runServerAddr string

type Config struct {
	Addr string `env:"ADDRESS"`
}

func initFlags() {
	pflag.StringVarP(&runServerAddr, "addr", "a", ":8080", "server address")
	pflag.Parse()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	if cfg.Addr != "" {
		runServerAddr = cfg.Addr
	}
}
