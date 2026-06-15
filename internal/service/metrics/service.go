package metrics

import (
	"slices"
	"strings"

	models "github.com/Flash0673/metrics-go/internal/model"
)

type Repository interface {
	Get(name, mType string) (*models.Metrics, error)
	Set(name, mType string, value any) (*models.Metrics, error)
	GetAll() []*models.Metrics
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Set(name, mType string, value any) (*models.Metrics, error) {
	return s.repo.Set(name, mType, value)
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
