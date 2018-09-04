package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
)

type (
	preferenceAllRequest struct {
		PhoneID string `json:"phoneId"`
	}

	preferenceAllResponse struct {
		Result []string
		Error  *APIError `json:"error,omitempty"`
	}

	preferenceUserAllEndpoint struct{}
)

var preferenceAll postEndpoint = preferenceUserAllEndpoint{}

func (r preferenceUserAllEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(preferenceAllRequest)

	logger := shared.GetLogger()
	logger = logger.With().
		Str("endpoint", r.GetPath()).
		Str("phoneId", request.PhoneID).Logger()
	logger.Info().Msgf("preference all request")

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	preferences, err := rtr.engines.PreferenceAll(request.PhoneID)

	var cuisines []string
	for _, preference := range preferences {
		cuisines = append(cuisines, preference.Cuisine)
	}

	result := preferenceAllResponse{Result: cuisines, Error: NewAPIError(err)}
	return result, err
}

func (r preferenceUserAllEndpoint) Validate(request interface{}) error {
	input := request.(preferenceAllRequest)
	if strings.TrimSpace(input.PhoneID) == "" {
		return helper.ValidationError{Message: fmt.Sprint("preference All failed, missing fields")}
	}

	return nil
}

func (r preferenceUserAllEndpoint) GetPath() string {
	return "/preference/all"
}

func (r preferenceUserAllEndpoint) HTTPRequest() interface{} {
	return preferenceAllRequest{}
}

func (r preferenceUserAllEndpoint) HTTPResult() interface{} {
	return preferenceAllResponse{}
}
