package http

import (
	"example/pkg/goruntime"
	"net/http"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer goruntime.RecoverPanic(r.Context(), "http-request-handling")
		next.ServeHTTP(w, r)
	})
}
