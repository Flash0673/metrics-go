package main

import (
	"net/http"

	"github.com/Flash0673/metrics-go/internal/handler"
	"github.com/Flash0673/metrics-go/internal/repository"
	"github.com/Flash0673/metrics-go/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	initFlags()

	repoAgg := repository.NewAggregator()
	svcAgg := service.NewAggregator(repoAgg)
	handlerAgg := handler.NewAggregator(svcAgg)

	mux := chi.NewRouter()
	mux.Get("/", handlerAgg.GetAll.ServeHTTP)
	mux.Get("/value/{type}/{name}", handlerAgg.Get.ServeHTTP)
	mux.Post("/update/{type}/{name}/{value}", handlerAgg.UpdateMetrics.ServeHTTP)
	if err := http.ListenAndServe(runServerAddr, mux); err != nil {
		panic(err)
	}
}
