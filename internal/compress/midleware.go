package compress

import (
	"net/http"
	"strings"
)

var contentTypesToBeEncoded = []string{"application/json", "text/html"}

func containsString(list []string, a string) bool {
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
		shouldBeEncoded := containsString(
			contentTypesToBeEncoded,
			w.Header().Get("Content-Type"),
		)
		if supportsGzip && shouldBeEncoded {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		if r.Header.Get("Content-Encoding") == "gzip" {
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
