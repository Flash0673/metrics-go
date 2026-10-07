package middleware

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressWriter_Middleware(t *testing.T) {
	t.Run("should write a compressed response", func(t *testing.T) {
		type s struct {
			A int `json:"a"`
			B int `json:"b"`
		}

		raw, _ := json.Marshal(&s{
			A: 1,
			B: 1,
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept-Encoding", "br, gzip")
		CompressWriterMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write(raw)
				assert.NoError(t, err)
			}),
		).ServeHTTP(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		gr, err := gzip.NewReader(resp.Body)
		assert.NoError(t, err)
		defer gr.Close()

		b := &bytes.Buffer{}
		_, err = b.ReadFrom(gr)
		assert.NoError(t, err)

		assert.Equal(t, raw, b.Bytes())
	})
}
