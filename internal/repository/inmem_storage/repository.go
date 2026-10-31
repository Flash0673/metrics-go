package inmem_storage

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	models "github.com/Flash0673/metrics-go/internal/model"
	"github.com/Flash0673/metrics-go/internal/repository/repo_err"
	"github.com/Flash0673/metrics-go/pkg/file"
	"github.com/Flash0673/metrics-go/pkg/logging/logger"
	"go.uber.org/zap"
)

type MemStorage struct {
	mu             *sync.RWMutex
	storage        map[string]*models.Metrics
	idGetter       func() string
	fStoreInterval time.Duration
	fStore         string
}

func NewMemStorage(ctx context.Context, restore bool, fStore string, fStoreInterval time.Duration) *MemStorage {
	s := make(map[string]*models.Metrics)

	if restore {
		data, err := file.Load(fStore)
		if err != nil {
			logger.Logger.Error("error loading file", zap.String("path", fStore))
		}
		err = json.Unmarshal(data, &s)
		if err != nil {
			logger.Logger.Error("error unmarshalling file", zap.String("path", fStore))
		}
	}

	ms := &MemStorage{
		mu:             new(sync.RWMutex),
		storage:        s,
		idGetter:       newIDGetter(),
		fStore:         fStore,
		fStoreInterval: fStoreInterval,
	}

	if fStoreInterval != time.Duration(0) {
		go func() {
			storeTicker := time.NewTicker(fStoreInterval)
			defer storeTicker.Stop()

		LOOP:
			for {
				select {
				case <-ctx.Done():
					break LOOP
				case <-storeTicker.C:
					ms.storeToFile()
				}
			}
		}()
	}

	return ms
}

func (m *MemStorage) GetAll() []*models.Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]*models.Metrics, 0, len(m.storage))
	for _, m := range m.storage {
		res = append(res, m)
	}
	return res
}

func (m *MemStorage) Get(name, mType string) (*models.Metrics, error) {
	key := generateKey(name, mType)
	m.mu.RLock()
	metrics, ok := m.storage[key]
	m.mu.RUnlock()
	if !ok {
		return &models.Metrics{
			ID:    name,
			MType: mType,
			Delta: nil,
			Value: nil,
			Hash:  key,
		}, repo_err.ErrNotFound
	}
	return metrics, nil
}

func (m *MemStorage) Set(name, mType string, metrics *models.Metrics) (*models.Metrics, error) {
	key := generateKey(name, mType)
	m.mu.Lock()
	m.storage[key] = metrics
	if m.fStoreInterval == time.Duration(0) {
		m.storeToFile()
	}
	m.mu.Unlock()
	return metrics, nil
}

func (m *MemStorage) storeToFile() {
	raw, err := json.Marshal(m.storage)
	if err != nil {
		logger.Logger.Error("error marshalling while saving to file")
	}
	err = file.Store(m.fStore, raw)
	if err != nil {
		logger.Logger.Error("error storing metrics to file", zap.String("error", err.Error()))
	}
}

func generateKey(name, mType string) string {
	return fmt.Sprintf("%s_%s", name, mType)
}

func newIDGetter() func() string {
	i := 0
	return func() string {
		i++
		return fmt.Sprintf("%d", i)
	}
}
