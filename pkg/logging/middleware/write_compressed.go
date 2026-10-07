package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/samber/lo"
)

type CompressWriter struct {
	w  http.ResponseWriter
	gr *gzip.Writer
}

func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	gr := gzip.NewWriter(w)
	return &CompressWriter{w: w, gr: gr}
}

func (c *CompressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *CompressWriter) Write(b []byte) (int, error) {
	return c.gr.Write(b)
}

func (c *CompressWriter) WriteHeader(statusCode int) {
	c.w.WriteHeader(statusCode)
}

func (c *CompressWriter) Close() error {
	return c.gr.Close()
}

func CompressWriterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncodings := r.Header.Get("Accept-Encoding")
		encs := strings.Split(acceptEncodings, ",")
		for idx, e := range encs {
			encs[idx] = strings.TrimSpace(e)
		}

		if lo.Contains(encs, "gzip") {
			cw := NewCompressWriter(w)
			ow = cw
			ow.Header().Set("Content-Encoding", "gzip")
			defer cw.Close()
		}

		next.ServeHTTP(ow, r)
	})
}
