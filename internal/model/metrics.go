package models

import (
	"fmt"
	"strings"
)

//go:generate easyjson -all metrics.go

const (
	Counter = "counter"
	Gauge   = "gauge"
)

// Metrics
// NOTE: Не усложняем пример, вводя иерархическую вложенность структур.
// Органичиваясь плоской моделью.
// Delta и Value объявлены через указатели,
// что бы отличать значение "0", от не заданного значения
// и соответственно не кодировать в структуру.
//
//easyjson:json
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"-"`
}

func (m *Metrics) GetName() string {
	return m.ID
}

func (m *Metrics) GetValue() string {
	switch m.MType {
	case Counter:
		return fmt.Sprintf("%d", *m.Delta)
	case Gauge:
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.10f", *m.Value), "0"), ".")
	}
	return ""
}

func (m *Metrics) GetDelta() int64 {
	if m.Delta == nil {
		return 0
	}
	return *m.Delta
}
