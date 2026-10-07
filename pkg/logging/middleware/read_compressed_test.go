package middleware

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressedReaderMiddleware(t *testing.T) {
	t.Run("should read a compressed request", func(t *testing.T) {
		type s struct {
			A int `json:"a"`
			B int `json:"b"`
		}

		raw, _ := json.Marshal(&s{
			A: 1,
			B: 1,
		})
		compressed := &bytes.Buffer{}
		gw := gzip.NewWriter(compressed)
		gw.Write(raw)
		gw.Close()

		h := CompressedReaderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			b, _ := io.ReadAll(req.Body)
			assert.Equal(t, raw, b)
		}))
		r := httptest.NewRequest("GET", "/", compressed)
		r.Header.Set("Content-Encoding", "gzip")
		h.ServeHTTP(httptest.NewRecorder(), r)
	})
}
