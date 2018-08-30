package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
)

type (
	preferenceDeleteRequest struct {
		PhoneID string   `json:"phoneId"`
		Cuisine []string `json:"cuisine,omitempty"`
	}

	preferenceDeleteResponse struct {
		preferenceDeleteRequest
		Error *APIError `json:"error,omitempty"`
	}

	preferenceUserDeleteEndpoint struct{}
)

var preferenceDelete postEndpoint = preferenceUserDeleteEndpoint{}

func (r preferenceUserDeleteEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(preferenceDeleteRequest)

	logger := shared.GetLogger()
	logger = logger.With().
		Str("endpoint", r.GetPath()).
		Str("phoneId", request.PhoneID).
		Strs("cuisine", request.Cuisine).Logger()
	logger.Info().Msgf("register phone request")

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	err := rtr.engines.PreferenceDelete(request.PhoneID, request.Cuisine)
	result := preferenceDeleteResponse{preferenceDeleteRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r preferenceUserDeleteEndpoint) Validate(request interface{}) error {
	input := request.(preferenceDeleteRequest)
	if strings.TrimSpace(input.PhoneID) == "" || len(input.Cuisine) == 0 {
		return helper.ValidationError{Message: fmt.Sprint("preference Delete failed, missing fields")}
	}

	return nil
}

func (r preferenceUserDeleteEndpoint) GetPath() string {
	return "/preference/delete"
}

func (r preferenceUserDeleteEndpoint) HTTPRequest() interface{} {
	return preferenceDeleteRequest{}
}

func (r preferenceUserDeleteEndpoint) HTTPResult() interface{} {
	return preferenceDeleteResponse{}
}
