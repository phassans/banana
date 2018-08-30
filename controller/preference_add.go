package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
)

type (
	preferenceAddRequest struct {
		PhoneID string   `json:"phoneId"`
		Cuisine []string `json:"cuisine,omitempty"`
	}

	preferenceAddResponse struct {
		preferenceAddRequest
		Error *APIError `json:"error,omitempty"`
	}

	preferenceUserAddEndpoint struct{}
)

var preferenceAdd postEndpoint = preferenceUserAddEndpoint{}

func (r preferenceUserAddEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(preferenceAddRequest)

	logger := shared.GetLogger()
	logger = logger.With().
		Str("endpoint", r.GetPath()).
		Str("phoneId", request.PhoneID).
		Strs("cuisine", request.Cuisine).Logger()
	logger.Info().Msgf("register phone request")

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	err := rtr.engines.PreferenceAdd(request.PhoneID, request.Cuisine)
	result := preferenceAddResponse{preferenceAddRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r preferenceUserAddEndpoint) Validate(request interface{}) error {
	input := request.(preferenceAddRequest)
	if strings.TrimSpace(input.PhoneID) == "" || len(input.Cuisine) == 0 {
		return helper.ValidationError{Message: fmt.Sprint("preference add failed, missing fields")}
	}

	return nil
}

func (r preferenceUserAddEndpoint) GetPath() string {
	return "/preference/add"
}

func (r preferenceUserAddEndpoint) HTTPRequest() interface{} {
	return preferenceAddRequest{}
}

func (r preferenceUserAddEndpoint) HTTPResult() interface{} {
	return preferenceAddResponse{}
}
