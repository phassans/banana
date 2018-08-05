package helper

import (
	"net"
	"net/http"

	"github.com/phassans/banana/shared"
)

// SetJSONContentResponse sets content type of response to be JSON.
func SetJSONContentResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next.ServeHTTP(w, r)
	})
}

func SetFieldsInLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := shared.GetLogger()
		ip, port, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
		}
		logger = logger.With().
			Str("ip", ip).
			Str("port", port).Logger()
		logger.Info().Msg("")
		next.ServeHTTP(w, r)
	})
}
