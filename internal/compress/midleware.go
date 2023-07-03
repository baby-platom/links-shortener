package compress

import (
	"net/http"
	"strings"

	"github.com/baby-platom/links-shortener/internal/logger"
)

var contentTypesToBeEncoded = []string{"application/json", "text/html"}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Middleware for compression and decompression
func Middleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		shouldBeEncoded := stringInSlice(
			r.Header.Get("Content-Type"),
			contentTypesToBeEncoded,
		)
		if supportsGzip && shouldBeEncoded {
			logger.Log.Info("Encoding content")
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	}
	return http.HandlerFunc(fn)
}
