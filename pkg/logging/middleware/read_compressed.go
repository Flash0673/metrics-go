package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
)

type CompressReader struct {
	r  io.ReadCloser
	gr *gzip.Reader
}

func (cr *CompressReader) Read(b []byte) (int, error) {
	return cr.gr.Read(b)
}

func (cr *CompressReader) Close() error {
	err := cr.r.Close()
	if err != nil {
		return err
	}

	return cr.gr.Close()
}

func NewCompressedReader(r io.ReadCloser) (*CompressReader, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &CompressReader{
		r:  r,
		gr: gr,
	}, nil
}

func CompressedReaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			cr, err := NewCompressedReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(w, r)
	})
}
