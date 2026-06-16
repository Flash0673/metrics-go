package update_metrics

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Flash0673/metrics-go/internal/handler/update_metrics/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUpdateMetrics(t *testing.T) {
	t.Parallel()

	type args struct {
		method      string
		metricsType string
		metricsName string
		metricValue string
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
				method:      http.MethodPost,
				metricsType: "gauge",
				metricsName: "test",
				metricValue: "1.1",
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		"counter success": {
			args: args{
				method:      http.MethodPost,
				metricsType: "counter",
				metricsName: "test",
				metricValue: "1",
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		"gauge error": {
			args: args{
				method:      http.MethodPost,
				metricsType: "gauge",
				metricsName: "test",
				metricValue: "s",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		"counter error": {
			args: args{
				method:      http.MethodPost,
				metricsType: "counter",
				metricsName: "test",
				metricValue: "s",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		"counter wrong value": {
			args: args{
				method:      http.MethodPost,
				metricsType: "counter",
				metricsName: "test",
				metricValue: "1.2",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		"wrong metric type": {
			args: args{
				method:      http.MethodPost,
				metricsType: "wrong",
				metricsName: "test",
				metricValue: "1.1",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		"empty metric name": {
			args: args{
				method:      http.MethodPost,
				metricsType: "gauge",
				metricsName: "",
				metricValue: "1",
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		"wrong method": {
			args: args{
				method:      http.MethodGet,
				metricsType: "gauge",
				metricsName: "test",
				metricValue: "1.2",
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
	r.Post("/update/{type}/{name}/{value}", h.ServeHTTP)
	s := httptest.NewServer(r)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			c := resty.New()
			resp, err := c.R().Execute(tc.args.method, fmt.Sprintf(
				"%s/update/%s/%s/%v",
				s.URL,
				tc.args.metricsType,
				tc.args.metricsName,
				tc.args.metricValue,
			))

			require.NoError(t, err)
			require.Equal(t, tc.want.statusCode, resp.StatusCode())
		})
	}
}
