package client

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Flash0673/metrics-go/internal/agent/dto"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    "http://localhost:8080",
	}
}

func (c *Client) ReportMetrics(metrics []dto.Metric) error {
	contentType := "text/plain"
	for _, m := range metrics {
		req, err := http.NewRequest(http.MethodPost, c.createUpdateUrl(m), nil)
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("Content-Type", contentType)
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return err
		} else {
			resp.Body.Close()
		}
	}
	return nil
}

func (c *Client) createUpdateUrl(m dto.Metric) string {
	parts := make([]string, 0, 5)
	parts = append(parts, c.baseURL)
	parts = append(parts, "update")
	parts = append(parts, m.GetType())
	parts = append(parts, m.GetName())
	parts = append(parts, m.GetValue())
	return strings.Join(parts, "/")
}
