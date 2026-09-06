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

	repoAgg := repository.NewAggregator()
	svcAgg := service.NewAggregator(repoAgg)
	handlerAgg := handler.NewAggregator(svcAgg)

	mux := chi.NewRouter()

	// Мидлеваре логирования реквеста
	mux.Use(middleware.RequestResponseLogging)

	mux.Get("/", handlerAgg.GetAll.ServeHTTP)
	mux.Get("/value/{type}/{name}", handlerAgg.Get.ServeHTTP)
	mux.Post("/update/{type}/{name}/{value}", handlerAgg.UpdateMetrics.ServeHTTP)
	if err := http.ListenAndServe(runServerAddr, mux); err != nil {
		panic(err)
	}
}
