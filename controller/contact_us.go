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
		PhoneID  string `json:"phoneId"`
		Name     string `json:"name,omitempty"`
		Email    string `json:"email,omitempty"`
		Subject  string `json:"subject,omitempty"`
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
		Str("phoneId", request.PhoneID).
		Str("name", request.Name).
		Str("subject", request.Subject).
		Str("email", request.Email).
		Str("comments", request.Comments).Logger()
	logger.Info().Msgf("contact us request")

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	err := rtr.engines.ContactUs(request.PhoneID, request.Name, request.Email, request.Subject, request.Comments)
	result := contactUsResponse{contactUsRequest: request, Error: NewAPIError(err)}
	return result, err
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
