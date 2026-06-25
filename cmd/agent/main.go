package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/Flash0673/metrics-go/internal/agent"
)

func main() {
	initFlags()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	ag := agent.New(targetUrl, reportInterval, pollInterval)

	ag.Run(ctx)

	<-ctx.Done()
	stop()
	// time to shout down
	fmt.Println("shouting down for 2 sec...")
	time.Sleep(2 * time.Second)
	fmt.Println("done")
}
