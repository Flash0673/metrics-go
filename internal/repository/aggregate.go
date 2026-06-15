package repository

import "github.com/Flash0673/metrics-go/internal/repository/inmem_storage"

type Aggregator struct {
	MemStorage *inmem_storage.MemStorage
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		MemStorage: inmem_storage.NewMemStorage(),
	}
}
