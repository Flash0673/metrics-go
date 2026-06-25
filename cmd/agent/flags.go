package main

import (
	"time"

	"github.com/spf13/pflag"
)

var (
	targetUrl      string
	reportInterval time.Duration
	pollInterval   time.Duration
)

func initFlags() {
	pflag.StringVarP(&targetUrl, "target-url", "t", "http://localhost:8080", "target url")
	pflag.DurationVarP(&reportInterval, "report-interval", "r", 10*time.Second, "report interval")
	pflag.DurationVarP(&pollInterval, "poll-interval", "p", 2*time.Second, "poll interval")

	pflag.Parse()
}
