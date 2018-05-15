package cachecontrol

import (
	"net/http"

	"github.com/acoshift/middleware"
)

// Config is cachecontrol config,
// empty string will be skipped
//
// Config[0] is fallback
type Config map[int]string

// New creates new cachecontrol middleware
func New(c Config) middleware.Middleware {
	return func(h http.Handler) http.Handler {
		if c == nil || len(c) == 0 {
			// by-pass middleware
			return h
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nw := &responseWriter{
				ResponseWriter: w,
				c:              c,
			}
			h.ServeHTTP(nw, r)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	c           Config
	wroteHeader bool
}

func (w *responseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true

	cc, ok := w.c[code]
	if !ok {
		cc = w.c[0]
	}
	if cc != "" {
		w.Header().Set("Cache-Control", cc)
	}

	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(p)
}
