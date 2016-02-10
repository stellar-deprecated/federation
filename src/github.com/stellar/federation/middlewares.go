package federation

import (
	"net/http"
)

func headersMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
