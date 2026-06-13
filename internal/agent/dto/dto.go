package dto

import "fmt"

type Metric interface {
	GetName() string
	GetType() string
	GetValue() string
}

type Gauge struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func NewGauge(name string, value float64) *Gauge {
	return &Gauge{
		Name:  name,
		Value: value,
	}
}

func (g *Gauge) GetName() string {
	return g.Name
}

func (g *Gauge) GetType() string {
	return "gauge"
}

func (g *Gauge) GetValue() string {
	return fmt.Sprintf("%f", g.Value)
}

type Counter struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

func NewCounter(name string, value int64) *Counter {
	return &Counter{
		Name:  name,
		Value: value,
	}
}

func (c *Counter) GetName() string {
	return c.Name
}

func (c *Counter) GetType() string {
	return "counter"
}

func (c *Counter) GetValue() string {
	return fmt.Sprintf("%d", c.Value)
}
