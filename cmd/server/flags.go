package main

import (
	"github.com/caarlos0/env/v11"
	"github.com/spf13/pflag"
)

var (
	runServerAddr string
	logLevel      string
)

// Config .
type Config struct {
	Addr     string `env:"ADDRESS"`
	LogLevel string `env:"LOG_LEVEL"`
}

func initFlags() {
	pflag.StringVarP(&runServerAddr, "addr", "a", ":8080", "server address")
	pflag.StringVarP(&logLevel, "loglevel", "l", "info", "log level")
	pflag.Parse()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	if cfg.Addr != "" {
		runServerAddr = cfg.Addr
	}

	if cfg.LogLevel != "" {
		logLevel = cfg.LogLevel
	}
}
