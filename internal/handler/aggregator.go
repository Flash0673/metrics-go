package handler

import "github.com/Flash0673/metrics-go/internal/handler/update_metrics"

type Aggregator struct {
	UpdateMetrics *update_metrics.Handler
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		UpdateMetrics: update_metrics.NewHandler(),
	}
}
