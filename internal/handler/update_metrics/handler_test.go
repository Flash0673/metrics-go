package update_metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
				metricValue: "",
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
				metricValue: "",
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
				metricValue: "",
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

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(
				tc.args.method,
				//fmt.Sprintf(
				//	"/update/%s/%s/%v",
				//	tc.args.metricsType,
				//	tc.args.metricsName,
				//	tc.args.metricValue,
				//),
				"/update/{type}/{name}/{value}",
				nil,
			)
			req.SetPathValue("type", tc.args.metricsType)
			req.SetPathValue("name", tc.args.metricsName)
			req.SetPathValue("value", tc.args.metricValue)
			rw := httptest.NewRecorder()
			h := NewHandler()
			h.ServeHTTP(rw, req)
			resp := rw.Result()

			require.Equal(t, tc.want.statusCode, resp.StatusCode)
		})
	}
}
