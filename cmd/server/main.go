package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Flash0673/metrics-go/internal/handler"
	"github.com/Flash0673/metrics-go/internal/repository"
	"github.com/Flash0673/metrics-go/internal/service"
	"github.com/Flash0673/metrics-go/pkg/logging/logger"
	"github.com/Flash0673/metrics-go/pkg/logging/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	initFlags()
	err := logger.New(logLevel)
	if err != nil {
		logger.Logger.Fatal("failed to init logger")
	}

	logger.Logger.Info("logger initialized")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	repoAgg := repository.NewAggregator(ctx, restore, fileStoragePath, storeInterval)
	svcAgg := service.NewAggregator(repoAgg)
	handlerAgg := handler.NewAggregator(svcAgg)

	mux := chi.NewRouter()

	// Мидлеваре логирования реквеста
	mux.Use(middleware.RequestResponseLogging)
	// Мидлеваре сжатия реквеста и респонса
	mux.Use(middleware.CompressedReaderMiddleware)
	mux.Use(middleware.CompressWriterMiddleware)

	mux.Get("/", handlerAgg.GetAll.ServeHTTP)
	mux.Post("/value/", handlerAgg.GetJSON.ServeHTTP)
	mux.Get("/value/{type}/{name}", handlerAgg.Get.ServeHTTP)
	mux.Post("/update/{type}/{name}/{value}", handlerAgg.UpdateMetrics.ServeHTTP)
	mux.Post("/update/", handlerAgg.UpdateMetricsJson.ServeHTTP)

	srv := &http.Server{
		Addr:    runServerAddr,
		Handler: mux,
	}

	// запускаем сервер в горутине
	go func() {
		logger.Logger.Info("server initialized")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Logger.Fatal("server failed", zap.Error(err))
		}
	}()

	// ждём сигнал
	<-ctx.Done()
	stop()

	logger.Logger.Info("shutting down...")

	// контекст с таймаутом для Shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Logger.Error("graceful shutdown failed", zap.Error(err))
	}
	logger.Logger.Info("done")
}
