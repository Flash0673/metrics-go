package service

import (
	"github.com/Flash0673/metrics-go/internal/repository"
	"github.com/Flash0673/metrics-go/internal/service/metrics"
)

type Aggregator struct {
	Metrics *metrics.Service
}

func NewAggregator(repoAgg *repository.Aggregator) *Aggregator {
	return &Aggregator{
		Metrics: metrics.NewService(repoAgg.MemStorage),
	}
}
