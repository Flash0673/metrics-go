package handler

import (
	"github.com/Flash0673/metrics-go/internal/handler/get"
	"github.com/Flash0673/metrics-go/internal/handler/get_all_metrics"
	"github.com/Flash0673/metrics-go/internal/handler/get_json"
	"github.com/Flash0673/metrics-go/internal/handler/update_metrics"
	"github.com/Flash0673/metrics-go/internal/handler/update_metrics_json"
	"github.com/Flash0673/metrics-go/internal/service"
)

type Aggregator struct {
	UpdateMetrics     *update_metrics.Handler
	UpdateMetricsJson *update_metrics_json.Handler
	Get               *get.Handler
	GetJSON           *get_json.Handler
	GetAll            *get_all_metrics.Handler
}

func NewAggregator(serviceAgg *service.Aggregator) *Aggregator {
	return &Aggregator{
		UpdateMetrics:     update_metrics.NewHandler(serviceAgg.Metrics),
		UpdateMetricsJson: update_metrics_json.NewHandler(serviceAgg.Metrics),
		Get:               get.NewHandler(serviceAgg.Metrics),
		GetJSON:           get_json.NewHandler(serviceAgg.Metrics),
		GetAll:            get_all_metrics.NewHandler(serviceAgg.Metrics),
	}
}
