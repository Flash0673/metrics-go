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
)

type Config struct {
	Addr           string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func initFlags() {
	pflag.StringVarP(&addr, "addr", "a", "localhost:8080", "target url")
	pflag.IntVarP(&reportIntervalInt, "report-interval", "r", 10, "report interval")
	pflag.IntVarP(&pollIntervalInt, "poll-interval", "p", 2, "poll interval")

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

	reportInterval = time.Duration(reportIntervalInt) * time.Second
	pollInterval = time.Duration(pollIntervalInt) * time.Second
}
