package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/phassans/banana/shared"
)

func (rtr *router) newPostHandler(endpoint postEndpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer rtr.cleanup(&err, w)

		logger := shared.GetLogger()

		request := reflect.New(reflect.TypeOf(endpoint.HTTPRequest()))
		err = json.NewDecoder(r.Body).Decode(request.Interface())
		if err != nil {
			//err = ErrInvalidJSON{Err: err}
			return
		}
		r.Body.Close()

		var result interface{}
		var engineErr error
		err = hystrixCall(endpoint, func() error {
			result, engineErr = endpoint.Execute(r.Context(), rtr, request.Elem().Interface())
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
			Str("query", fmt.Sprintf("%#v", request.Elem().Interface())).
			Int("errorStatus", GetErrorStatus(err)).Logger()
		if err != nil {
			logger.Info().Msgf("POST failure")
			return
		}
		logger.Info().Msgf("POST success")

		err = json.NewEncoder(w).Encode(result)
	}
}
