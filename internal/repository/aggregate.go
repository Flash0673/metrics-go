package repository

import (
	"context"
	"time"

	"github.com/Flash0673/metrics-go/internal/repository/db"
	"github.com/Flash0673/metrics-go/internal/repository/inmem_storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Aggregator struct {
	MemStorage *inmem_storage.MemStorage
	DB         *db.Storage
}

func NewAggregator(ctx context.Context, dbpool *pgxpool.Pool, restore bool, fStore string, fStoreInterval time.Duration) *Aggregator {
	return &Aggregator{
		MemStorage: inmem_storage.NewMemStorage(ctx, restore, fStore, fStoreInterval),
		DB:         db.New(ctx, dbpool),
	}
}
