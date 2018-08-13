package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/banana/helper"
)

type (
	registerPhoneRequest struct {
		RegistrationToken string `json:"registrationToken"`
		PhoneID           string `json:"phoneId"`
		PhoneModel        string `json:"phoneModel"`
	}

	registerPhoneResponse struct {
		registerPhoneRequest
		Error *APIError `json:"error,omitempty"`
	}

	registerPhoneEndpoint struct{}
)

var registerPhone postEndpoint = registerPhoneEndpoint{}

func (r registerPhoneEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(registerPhoneRequest)
	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	err := rtr.engines.RegisterPhone(request.RegistrationToken, request.PhoneID, request.PhoneModel)
	result := registerPhoneResponse{registerPhoneRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r registerPhoneEndpoint) Validate(request interface{}) error {
	input := request.(registerPhoneRequest)
	if strings.TrimSpace(input.RegistrationToken) == "" ||
		strings.TrimSpace(input.PhoneID) == "" {
		return helper.ValidationError{Message: fmt.Sprint("phone registration failed, missing fields")}
	}

	return nil
}

func (r registerPhoneEndpoint) GetPath() string {
	return "/register/phone"
}

func (r registerPhoneEndpoint) HTTPRequest() interface{} {
	return registerPhoneRequest{}
}

func (r registerPhoneEndpoint) HTTPResult() interface{} {
	return registerPhoneResponse{}
}
