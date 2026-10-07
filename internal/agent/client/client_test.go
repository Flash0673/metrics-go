package client

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Flash0673/metrics-go/internal/agent/dto"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
)

func TestClient_reportMetricsJson(t *testing.T) {
	t.Run("compression test", func(t *testing.T) {
		m := dto.NewCounter("test", 1)
		mJson, err := easyjson.Marshal(m.ToModel())
		assert.NoError(t, err)

		h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ce := req.Header.Get("Content-Encoding")
			assert.Equal(t, "gzip", ce)

			defer req.Body.Close()
			r, err := gzip.NewReader(req.Body)
			defer r.Close()
			assert.NoError(t, err)

			res := bytes.Buffer{}
			_, err = res.ReadFrom(r)
			assert.NoError(t, err)

			assert.True(t, bytes.Equal(res.Bytes(), mJson))

			w.WriteHeader(http.StatusOK)
		})

		srv := httptest.NewServer(h)
		defer srv.Close()

		c := &Client{
			httpClient: resty.New(),
			baseURL:    srv.URL,
		}
		err = c.reportMetricsJson([]dto.Metric{m})
		assert.NoError(t, err)
	})
}
