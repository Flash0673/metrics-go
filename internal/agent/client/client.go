package client

import (
	"fmt"

	"github.com/Flash0673/metrics-go/internal/agent/dto"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	httpClient *resty.Client
	baseURL    string
}

func NewClient(targetUrl string) *Client {
	return &Client{
		httpClient: resty.New(),
		baseURL:    targetUrl,
	}
}

func (c *Client) ReportMetrics(metrics []dto.Metric) error {
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
