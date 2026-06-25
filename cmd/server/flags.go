package main

import "github.com/spf13/pflag"

var runServerAddr string

func initFlags() {
	pflag.StringVarP(&runServerAddr, "addr", "a", ":8080", "server address")
	pflag.Parse()
}
