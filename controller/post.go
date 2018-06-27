package controller

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/rs/xlog"
)

func (rtr *router) newPostHandler(endpoint postEndpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer rtr.cleanup(&err, w)

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
		if err != nil {
			xlog.Warnf("POST %s query %+v error %d: %s", endpoint.GetPath(), request.Elem().Interface(), GetErrorStatus(err), err.Error())
			return
		}

		err = json.NewEncoder(w).Encode(result)
	}
}
