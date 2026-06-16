package inmem_storage

import (
	"fmt"
	"sync"

	models "github.com/Flash0673/metrics-go/internal/model"
	"github.com/Flash0673/metrics-go/internal/repository/repo_err"
)

type MemStorage struct {
	mu       *sync.RWMutex
	storage  map[string]*models.Metrics
	idGetter func() string
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		mu:       new(sync.RWMutex),
		storage:  make(map[string]*models.Metrics),
		idGetter: newIDGetter(),
	}
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
			ID:    m.idGetter(),
			MType: mType,
			Delta: nil,
			Value: nil,
			Hash:  key,
		}, repo_err.ErrNotFound
	}
	return metrics, nil
}

func (m *MemStorage) Set(name, mType string, metrics *models.Metrics) error {
	key := generateKey(name, mType)
	m.mu.Lock()
	m.storage[key] = metrics
	m.mu.Unlock()
	return nil
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
