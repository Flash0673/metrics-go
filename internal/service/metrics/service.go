package metrics

import (
	"errors"
	"slices"
	"strings"

	models "github.com/Flash0673/metrics-go/internal/model"
	"github.com/Flash0673/metrics-go/internal/repository/repo_err"
)

type Repository interface {
	Get(name, mType string) (*models.Metrics, error)
	Set(name, mType string, metrics *models.Metrics) error
	GetAll() []*models.Metrics
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Set(name, mType string, value any) error {
	metrics, err := s.repo.Get(name, mType)
	if err != nil && !errors.Is(err, repo_err.ErrNotFound) {
		return err
	}

	switch mType {
	case models.Gauge:
		v := value.(float64)
		metrics.Value = &v
	case models.Counter:
		newValue := metrics.GetDelta() + value.(int64)
		metrics.Delta = &newValue
	default:
		return errors.New("unsupported metric type")
	}

	return s.repo.Set(name, mType, metrics)
}

func (s *Service) Get(name, mType string) (*models.Metrics, error) {
	return s.repo.Get(name, mType)
}

func (s *Service) GetAll() []*models.Metrics {
	metrics := s.repo.GetAll()
	slices.SortFunc(metrics, func(a, b *models.Metrics) int {
		return strings.Compare(a.ID, b.ID)
	})
	return metrics
}
