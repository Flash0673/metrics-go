package main

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/spf13/pflag"
)

var (
	runServerAddr string
	logLevel      string

	// параметры сохранения в файл
	storeInterval   time.Duration
	fileStoragePath string
	restore         bool
)

// Config .
type Config struct {
	Addr     string `env:"ADDRESS"`
	LogLevel string `env:"LOG_LEVEL"`

	StoreInterval   time.Duration `env:"STORE_INTERVAL"`
	FileStoragePath string        `env:"FILE_STORAGE_PATH"`
	Restore         bool          `env:"RESTORE"`
}

func initFlags() {
	pflag.StringVarP(&runServerAddr, "addr", "a", ":8080", "server address")
	pflag.StringVarP(&logLevel, "loglevel", "l", "info", "log level")
	pflag.DurationVarP(&storeInterval, "store-interval", "i", 300*time.Second, "store interval")
	pflag.StringVarP(&fileStoragePath, "file-storage-path", "f", "/tmp/metrics-db.json", "file storage path")
	pflag.BoolVarP(&restore, "restore", "r", true, "restore")
	pflag.Parse()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	if cfg.Addr != "" {
		runServerAddr = cfg.Addr
	}

	if cfg.LogLevel != "" {
		logLevel = cfg.LogLevel
	}

	if cfg.StoreInterval != time.Duration(0) {
		storeInterval = cfg.StoreInterval
	}

	if cfg.FileStoragePath != "" {
		fileStoragePath = cfg.FileStoragePath
	}

	if cfg.Restore {
		restore = true
	}
}
