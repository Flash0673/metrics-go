package main

import (
	"time"

	"github.com/spf13/pflag"
)

var (
	addr              string
	reportIntervalInt int
	pollIntervalInt   int
	reportInterval    time.Duration
	pollInterval      time.Duration
)

func initFlags() {
	pflag.StringVarP(&addr, "addr", "a", "localhost:8080", "target url")
	pflag.IntVarP(&reportIntervalInt, "report-interval", "r", 10, "report interval")
	pflag.IntVarP(&pollIntervalInt, "poll-interval", "p", 2, "poll interval")

	pflag.Parse()

	reportInterval = time.Duration(reportIntervalInt) * time.Second
	pollInterval = time.Duration(pollIntervalInt) * time.Second
}
