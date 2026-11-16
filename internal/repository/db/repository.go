package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(_ context.Context, pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

func (s *Storage) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}
