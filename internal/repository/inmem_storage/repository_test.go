package inmem_storage

import (
	"encoding/json"
	"os"
	"testing"

	models "github.com/Flash0673/metrics-go/internal/model"
	"github.com/Flash0673/metrics-go/pkg/file"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestNewMemStorage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		storedMetrics := map[string]*models.Metrics{
			"first": {
				ID:    "1",
				MType: "gauge",
				Value: lo.ToPtr(20.0),
			},
			"second": {
				ID:    "2",
				MType: "counter",
				Delta: lo.ToPtr(int64(11)),
			},
		}

		tmpDir := "./tmp_metrics"
		defer os.Remove(tmpDir)

		prepareTmpFile(tmpDir, storedMetrics)

		s := NewMemStorage(true, tmpDir)
		require.Equal(t, storedMetrics, s.storage)
	})
}

func prepareTmpFile(source string, data map[string]*models.Metrics) {
	raw, _ := json.Marshal(data)
	file.Store(source, raw)
}
