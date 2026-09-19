package main

import (
	"net/http"

	"github.com/Flash0673/metrics-go/internal/handler"
	"github.com/Flash0673/metrics-go/internal/repository"
	"github.com/Flash0673/metrics-go/internal/service"
	"github.com/Flash0673/metrics-go/pkg/logging/logger"
	"github.com/Flash0673/metrics-go/pkg/logging/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	initFlags()
	err := logger.New(logLevel)
	if err != nil {
		logger.Logger.Fatal("failed to init logger")
	}

	logger.Logger.Info("logger initialized")

	repoAgg := repository.NewAggregator()
	svcAgg := service.NewAggregator(repoAgg)
	handlerAgg := handler.NewAggregator(svcAgg)

	mux := chi.NewRouter()

	// Мидлеваре логирования реквеста
	mux.Use(middleware.RequestResponseLogging)

	mux.Get("/", handlerAgg.GetAll.ServeHTTP)
	mux.Post("/value", handlerAgg.GetJSON.ServeHTTP)
	mux.Get("/value/{type}/{name}", handlerAgg.Get.ServeHTTP)
	mux.Post("/update/{type}/{name}/{value}", handlerAgg.UpdateMetrics.ServeHTTP)
	mux.Post("/update", handlerAgg.UpdateMetricsJson.ServeHTTP)

	logger.Logger.Info("server initialized")
	if err := http.ListenAndServe(runServerAddr, mux); err != nil {
		panic(err)
	}
}
