package update_metrics_json

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Flash0673/metrics-go/internal/handler/update_metrics_json/mocks"
	models "github.com/Flash0673/metrics-go/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUpdateMetrics(t *testing.T) {
	t.Parallel()

	type args struct {
		method        string
		bodyGenerator func() string
	}

	type want struct {
		statusCode int
	}

	tests := map[string]struct {
		args args
		want want
	}{
		"gauge success": {
			args: args{
				method: http.MethodPost,
				bodyGenerator: func() string {
					v := 1.1
					m := models.Metrics{
						ID:    "test",
						MType: "gauge",
						Value: &v,
					}
					body, _ := json.Marshal(m)
					return string(body)
				},
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		"counter success": {
			args: args{
				method: http.MethodPost,
				bodyGenerator: func() string {
					v := int64(1)
					m := models.Metrics{
						ID:    "test",
						MType: "counter",
						Delta: &v,
					}
					body, _ := json.Marshal(m)
					return string(body)
				},
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		"gauge error": {
			args: args{
				method: http.MethodPost,
				bodyGenerator: func() string {
					m := models.Metrics{
						ID:    "test",
						MType: "gauge",
						Value: nil,
					}
					body, _ := json.Marshal(m)
					return string(body)
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		"counter error": {
			args: args{
				method: http.MethodPost,
				bodyGenerator: func() string {
					m := models.Metrics{
						ID:    "test",
						MType: "counter",
						Delta: nil,
					}
					body, _ := json.Marshal(m)
					return string(body)
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		"empty metric name": {
			args: args{
				method: http.MethodPost,
				bodyGenerator: func() string {
					v := 1.1
					m := models.Metrics{
						ID:    "",
						MType: "gauge",
						Value: &v,
					}
					body, _ := json.Marshal(m)
					return string(body)
				},
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		"wrong method": {
			args: args{
				method: http.MethodGet,
				bodyGenerator: func() string {
					v := 1.2
					m := models.Metrics{
						ID:    "test",
						MType: "gauge",
						Value: &v,
					}
					body, _ := json.Marshal(m)
					return string(body)
				},
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
	}
	ctrl := gomock.NewController(t)
	m := mocks.NewMockService(ctrl)
	m.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	h := NewHandler(m)
	r := chi.NewRouter()
	r.Post("/update", h.ServeHTTP)
	s := httptest.NewServer(r)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			c := resty.New()

			resp, err := c.R().SetBody(tc.args.bodyGenerator()).Execute(tc.args.method, fmt.Sprintf(
				"%s/update",
				s.URL,
			))

			require.NoError(t, err)
			require.Equal(t, tc.want.statusCode, resp.StatusCode())
		})
	}
}
