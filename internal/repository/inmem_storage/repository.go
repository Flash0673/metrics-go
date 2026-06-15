package inmem_storage

import (
	"errors"
	"fmt"

	models "github.com/Flash0673/metrics-go/internal/model"
)

var ErrNotFound error = errors.New("not found")

type MemStorage struct {
	storage  map[string]*models.Metrics
	idGetter func() string
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		storage:  make(map[string]*models.Metrics),
		idGetter: newIDGetter(),
	}
}

func (m *MemStorage) GetAll() []*models.Metrics {
	res := make([]*models.Metrics, 0, len(m.storage))
	for _, m := range m.storage {
		res = append(res, m)
	}
	return res
}

func (m *MemStorage) Get(name, mType string) (*models.Metrics, error) {
	key := generateKey(name, mType)
	metrics, ok := m.storage[key]
	if !ok {
		return nil, ErrNotFound
	}
	return metrics, nil
}

func (m *MemStorage) Set(name, mType string, value any) (*models.Metrics, error) {
	key := generateKey(name, mType)
	metrics, ok := m.storage[key]
	if !ok {
		metrics = &models.Metrics{
			ID:    m.idGetter(),
			MType: mType,
			Delta: nil,
			Value: nil,
			Hash:  key,
		}
	}
	switch mType {
	case models.Counter:
		v := value.(int64)
		metrics.Delta = &v
	case models.Gauge:
		v := value.(float64)
		metrics.Value = &v
	default:
		return nil, fmt.Errorf("unknown metrics type: %s", mType)
	}
	m.storage[key] = metrics
	return metrics, nil
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
