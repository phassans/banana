package helper

import (
	"net/http"
)

// SetJSONContentResponse sets content type of response to be JSON.
func SetJSONContentResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next.ServeHTTP(w, r)
	})
}
