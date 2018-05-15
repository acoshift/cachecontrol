package cachecontrol_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/acoshift/cachecontrol"
	"github.com/stretchr/testify/assert"
)

func TestCacheControl(t *testing.T) {
	cc := cachecontrol.New(cachecontrol.Config{
		http.StatusOK:       "public, max-age=3600",
		http.StatusNotFound: "private, max-age=0",
		0:                   "private",
	})

	t.Run("OK", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		cc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(w, r)
		assert.Equal(t, "public, max-age=3600", w.Header().Get("Cache-Control"))
	})

	t.Run("OKWithoutWriteHeader", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		cc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})).ServeHTTP(w, r)
		assert.Equal(t, "public, max-age=3600", w.Header().Get("Cache-Control"))
	})

	t.Run("NotFound", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		cc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})).ServeHTTP(w, r)
		assert.Equal(t, "private, max-age=0", w.Header().Get("Cache-Control"))
	})

	t.Run("Fallback", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		cc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})).ServeHTTP(w, r)
		assert.Equal(t, "private", w.Header().Get("Cache-Control"))
	})

	t.Run("DoubleWriteHeaders", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		cc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.WriteHeader(http.StatusNotFound)
		})).ServeHTTP(w, r)
		assert.Equal(t, "public, max-age=3600", w.Header().Get("Cache-Control"))
	})
}

func TestBypass(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	assert.Equal(t, reflect.ValueOf(h).Pointer(), reflect.ValueOf(cachecontrol.New(nil)(h)).Pointer())
	assert.Equal(t, reflect.ValueOf(h).Pointer(), reflect.ValueOf(cachecontrol.New(cachecontrol.Config{})(h)).Pointer())
}
