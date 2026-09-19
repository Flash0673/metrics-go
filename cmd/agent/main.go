package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/Flash0673/metrics-go/internal/agent"
	"github.com/Flash0673/metrics-go/pkg/logging/logger"
)

func main() {
	initFlags()
	err := logger.New(logLevel)
	if err != nil {
		logger.Logger.Fatal("failed to init logger")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	logger.Logger.Info("starting agent")
	ag := agent.New(addr, reportInterval, pollInterval)

	ag.Run(ctx)

	<-ctx.Done()
	stop()
	// time to shout down
	fmt.Println("shouting down for 2 sec...")
	time.Sleep(2 * time.Second)
	fmt.Println("done")
}
