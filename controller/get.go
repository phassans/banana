package controller

import (
	"encoding/json"
	"net/http"

	"github.com/rs/xlog"
)

func (rtr *router) newGetHandler(endpoint getEndPoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer rtr.cleanup(&err, w)
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
		if err != nil {
			xlog.Warnf("GET %s query %+v error %d: %s", endpoint.GetPath(), r.URL.RawQuery, GetErrorStatus(err), err.Error())
			return
		} else {
			xlog.Infof("GET %s query %+v success %d", endpoint.GetPath(), r.URL.RawQuery, GetErrorStatus(err))
		}

		err = json.NewEncoder(w).Encode(result)
	}
}
