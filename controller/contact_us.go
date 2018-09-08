package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
)

type (
	contactUsRequest struct {
		Name     string `json:"name,omitempty"`
		Email    string `json:"email,omitempty"`
		Comments string `json:"comments"`
	}

	contactUsResponse struct {
		contactUsRequest
		Error *APIError `json:"error,omitempty"`
	}

	contactUsEndpoint struct{}
)

var contactUs postEndpoint = contactUsEndpoint{}

func (r contactUsEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(contactUsRequest)

	logger := shared.GetLogger()
	logger = logger.With().
		Str("endpoint", r.GetPath()).
		Str("name", request.Name).
		Str("email", request.Email).
		Str("comments", request.Comments).Logger()
	logger.Info().Msgf("contact us request")

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	//err := rtr.engines.RegisterPhone(request.RegistrationToken, request.PhoneID, request.PhoneModel)
	result := contactUsResponse{contactUsRequest: request, Error: NewAPIError(nil)}
	return result, nil
}

func (r contactUsEndpoint) Validate(request interface{}) error {
	input := request.(contactUsRequest)
	if strings.TrimSpace(input.Comments) == "" {
		return helper.ValidationError{Message: fmt.Sprint("contact us failed, missing fields")}
	}

	return nil
}

func (r contactUsEndpoint) GetPath() string {
	return "/contactus"
}

func (r contactUsEndpoint) HTTPRequest() interface{} {
	return contactUsRequest{}
}

func (r contactUsEndpoint) HTTPResult() interface{} {
	return contactUsResponse{}
}
