package agent

import (
	"context"
	"log"
	"math/rand"
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
	pollCounter                  int
}

func New(targetUrl string, reportInterval, pollInterval time.Duration) *Agent {
	return &Agent{
		rw: &sync.RWMutex{},
		// TODO add config
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		client:         client.NewClient(targetUrl),
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
	a.pollCounter++
}

func (a *Agent) reportMetrics() {
	a.rw.RLock()
	snapShot := a.m
	a.rw.RUnlock()

	metrcis := make([]dto.Metric, 0, 29)

	// TODO add another metrics
	metrcis = append(metrcis, dto.NewGauge("Alloc", float64(snapShot.Alloc)))
	metrcis = append(metrcis, dto.NewGauge("BuckHashSys", float64(snapShot.BuckHashSys)))
	metrcis = append(metrcis, dto.NewGauge("Frees", float64(snapShot.Frees)))
	metrcis = append(metrcis, dto.NewGauge("GCCPUFraction", snapShot.GCCPUFraction))
	metrcis = append(metrcis, dto.NewGauge("GCSys", float64(snapShot.GCSys)))
	metrcis = append(metrcis, dto.NewGauge("HeapAlloc", float64(snapShot.HeapAlloc)))
	metrcis = append(metrcis, dto.NewGauge("HeapIdle", float64(snapShot.HeapIdle)))
	metrcis = append(metrcis, dto.NewGauge("HeapInuse", float64(snapShot.HeapInuse)))
	metrcis = append(metrcis, dto.NewGauge("HeapObjects", float64(snapShot.HeapObjects)))
	metrcis = append(metrcis, dto.NewGauge("HeapReleased", float64(snapShot.HeapReleased)))
	metrcis = append(metrcis, dto.NewGauge("HeapSys", float64(snapShot.HeapSys)))
	metrcis = append(metrcis, dto.NewGauge("LastGC", float64(snapShot.LastGC)))
	metrcis = append(metrcis, dto.NewGauge("Lookups", float64(snapShot.Lookups)))
	metrcis = append(metrcis, dto.NewGauge("MCacheInuse", float64(snapShot.MCacheInuse)))
	metrcis = append(metrcis, dto.NewGauge("MCacheSys", float64(snapShot.MCacheSys)))
	metrcis = append(metrcis, dto.NewGauge("MSpanInuse", float64(snapShot.MSpanInuse)))
	metrcis = append(metrcis, dto.NewGauge("MSpanSys", float64(snapShot.MSpanSys)))
	metrcis = append(metrcis, dto.NewGauge("Mallocs", float64(snapShot.Mallocs)))
	metrcis = append(metrcis, dto.NewGauge("NextGC", float64(snapShot.NextGC)))
	metrcis = append(metrcis, dto.NewGauge("NumForcedGC", float64(snapShot.NumForcedGC)))
	metrcis = append(metrcis, dto.NewGauge("NumGC", float64(snapShot.NumGC)))
	metrcis = append(metrcis, dto.NewGauge("OtherSys", float64(snapShot.OtherSys)))
	metrcis = append(metrcis, dto.NewGauge("PauseTotalNs", float64(snapShot.PauseTotalNs)))
	metrcis = append(metrcis, dto.NewGauge("StackInuse", float64(snapShot.StackInuse)))
	metrcis = append(metrcis, dto.NewGauge("StackSys", float64(snapShot.StackSys)))
	metrcis = append(metrcis, dto.NewGauge("Sys", float64(snapShot.Sys)))
	metrcis = append(metrcis, dto.NewGauge("TotalAlloc", float64(snapShot.TotalAlloc)))

	metrcis = append(metrcis, dto.NewCounter("PollCount", int64(a.pollCounter)))
	metrcis = append(metrcis, dto.NewGauge("RandomValue", rand.Float64()*100))

	err := a.client.ReportMetrics(metrcis)
	if err != nil {
		log.Printf("Error reporting metrics: %v", err)
	}
}
