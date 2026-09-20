package client

import (
	"fmt"

	"github.com/Flash0673/metrics-go/internal/agent/dto"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
)

type ReportStrategy int

const (
	PathParams ReportStrategy = iota
	Body
)

type Client struct {
	httpClient *resty.Client
	baseURL    string
	strategy   ReportStrategy
}

func NewClient(addr string) *Client {
	return &Client{
		httpClient: resty.New(),
		baseURL:    fmt.Sprintf("http://%s", addr),
		strategy:   PathParams,
	}
}

func (c *Client) WithReportStrategy(s ReportStrategy) *Client {
	c.strategy = s
	return c
}

func (c *Client) SetReportStrategy(s ReportStrategy) {
	c.strategy = s
}

func (c *Client) ReportMetrics(metrics []dto.Metric) error {
	switch c.strategy {
	case PathParams:
		return c.reportMetrics(metrics)
	case Body:
		return c.reportMetricsJson(metrics)
	default:
		return c.reportMetricsJson(metrics)
	}
}

func (c *Client) reportMetrics(metrics []dto.Metric) error {
	contentType := "text/plain"
	for _, m := range metrics {
		_, err := c.httpClient.R().
			SetHeader("Content-Type", contentType).
			SetPathParams(map[string]string{
				"type":  m.GetType(),
				"name":  m.GetName(),
				"value": m.GetValue(),
			}).
			Post(fmt.Sprintf("%s/update/{type}/{name}/{value}", c.baseURL))
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) reportMetricsJson(metrics []dto.Metric) error {
	contentType := "application/json"
	for _, m := range metrics {
		body, err := easyjson.Marshal(m.ToModel())
		if err != nil {
			return err
		}
		_, err = c.httpClient.R().
			SetHeader("Content-Type", contentType).
			SetBody(body).
			Post(fmt.Sprintf("%s/update/", c.baseURL))
		if err != nil {
			return err
		}
	}
	return nil
}
