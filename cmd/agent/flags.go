package main

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/spf13/pflag"
)

var (
	addr              string
	reportIntervalInt int
	pollIntervalInt   int
	reportInterval    time.Duration
	pollInterval      time.Duration
	logLevel          string
)

type Config struct {
	Addr           string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	LogLevel       string `env:"LOG_LEVEL"`
}

func initFlags() {
	pflag.StringVarP(&addr, "addr", "a", "localhost:8080", "target url")
	pflag.IntVarP(&reportIntervalInt, "report-interval", "r", 10, "report interval")
	pflag.IntVarP(&pollIntervalInt, "poll-interval", "p", 2, "poll interval")
	pflag.StringVarP(&logLevel, "loglevel", "l", "info", "log level")

	pflag.Parse()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	if cfg.Addr != "" {
		addr = cfg.Addr
	}
	if cfg.ReportInterval != 0 {
		reportIntervalInt = cfg.ReportInterval
	}
	if cfg.PollInterval != 0 {
		pollIntervalInt = cfg.PollInterval
	}
	if cfg.LogLevel != "" {
		logLevel = cfg.LogLevel
	}

	reportInterval = time.Duration(reportIntervalInt) * time.Second
	pollInterval = time.Duration(pollIntervalInt) * time.Second
}
