package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/phassans/banana/shared"
)

func (rtr *router) newGetHandler(endpoint getEndPoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer rtr.cleanup(&err, w)

		logger := shared.GetLogger()

		vals := r.URL.Query()

		var result interface{}
		var engineErr error
		err = hystrixCall(endpoint, func() error {
			result, engineErr = endpoint.Do(r.Context(), rtr, vals)
			if IsHardError(engineErr) {
				return engineErr
			}
			return nil
		})
		if err == nil {
			err = engineErr
		}
		logger = logger.With().
			Str("endpoint", endpoint.GetPath()).
			Str("query", fmt.Sprintf("%#v", r.URL.RawQuery)).
			Int("errorStatus", GetErrorStatus(err)).Logger()

		if err != nil {
			logger.Error().Msgf(err.Error())
			return
		}
		logger.Info().Msgf("GET success")

		err = json.NewEncoder(w).Encode(result)
	}
}
