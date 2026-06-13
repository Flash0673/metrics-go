package agent

import (
	"context"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/Flash0673/metrics-go/internal/agent/client"
	"github.com/Flash0673/metrics-go/internal/agent/dto"
)

type Agent struct {
	m                            runtime.MemStats
	rw                           *sync.RWMutex
	pollInterval, reportInterval time.Duration
	client                       *client.Client
}

func New() *Agent {
	return &Agent{
		rw: &sync.RWMutex{},
		// TODO add config
		pollInterval:   2 * time.Second,
		reportInterval: 10 * time.Second,
		client:         client.NewClient(),
	}
}

func (a *Agent) Run(ctx context.Context) {
	go a.run(ctx)
}

func (a *Agent) run(ctx context.Context) {
	pollTicker := time.NewTicker(a.pollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(a.reportInterval)
	defer reportTicker.Stop()
LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case <-pollTicker.C:
			a.pollMetrics()
		case <-reportTicker.C:
			a.reportMetrics()
		default:
		}
	}
}

func (a *Agent) pollMetrics() {
	a.rw.Lock()
	defer a.rw.Unlock()
	runtime.ReadMemStats(&a.m)
}

func (a *Agent) reportMetrics() {
	a.rw.RLock()
	snapShot := a.m
	a.rw.RUnlock()

	metrcis := make([]dto.Metric, 0, 29)

	// TODO add another metrics
	metrcis = append(metrcis, dto.NewGauge("Alloc", float64(snapShot.Alloc)))

	err := a.client.ReportMetrics(metrcis)
	if err != nil {
		log.Printf("Error reporting metrics: %v", err)
	}
}
