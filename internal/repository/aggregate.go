package repository

import (
	"context"
	"time"

	"github.com/Flash0673/metrics-go/internal/repository/inmem_storage"
)

type Aggregator struct {
	MemStorage *inmem_storage.MemStorage
}

func NewAggregator(ctx context.Context, restore bool, fStore string, fStoreInterval time.Duration) *Aggregator {
	return &Aggregator{
		MemStorage: inmem_storage.NewMemStorage(ctx, restore, fStore, fStoreInterval),
	}
}
