package main

import (
	"net/http"

	"github.com/Flash0673/metrics-go/internal/handler"
)

func main() {
	handlerAgg := handler.NewAggregator()

	mux := http.NewServeMux()
	mux.Handle("/update/{type}/{name}/{value}", handlerAgg.UpdateMetrics)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
