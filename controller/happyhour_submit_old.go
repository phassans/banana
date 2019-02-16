package controller

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
)

type (
	happyHourSubmitRequest struct {
		PhoneID       string                  `json:"phoneId"`
		Name          string                  `json:"name,omitempty"`
		Email         string                  `json:"email,omitempty"`
		BusinessOwner bool                    `json:"businessOwner"`
		Restaurant    string                  `json:"restaurant,omitempty"`
		City          string                  `json:"city,omitempty"`
		Description   string                  `json:"description,omitempty"`
		images        []*multipart.FileHeader `json:"images,omitempty"`
	}

	happyHourSubmitResponse struct {
		happyHourSubmitRequest
		Error *APIError `json:"error,omitempty"`
	}

	happyHourSubmitEndpoint struct{}
)

var happyHourSubmit postEndpoint = happyHourSubmitEndpoint{}

func (r happyHourSubmitEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(happyHourSubmitRequest)

	logger := shared.GetLogger()
	logger = logger.With().
		Str("endpoint", r.GetPath()).
		Str("phoneId", request.PhoneID).
		Str("name", request.Name).
		Str("email", request.Email).
		Bool("BusinessOwner", request.BusinessOwner).
		Str("restaurant", request.Restaurant).
		Str("city", request.City).
		Str("description", request.Description).Logger()
	logger.Info().Msgf("submit happy hour request")

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	_, err := rtr.engines.SubmitHappyHour(request.PhoneID, request.Name, request.Email, request.BusinessOwner, request.Restaurant, request.City, request.Description)
	result := happyHourSubmitResponse{happyHourSubmitRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r happyHourSubmitEndpoint) Validate(request interface{}) error {
	input := request.(happyHourSubmitRequest)
	if strings.TrimSpace(input.Restaurant) == "" {
		return helper.ValidationError{Message: fmt.Sprint("submit happy hour failed, missing restaurant")}
	}
	if strings.TrimSpace(input.City) == "" {
		return helper.ValidationError{Message: fmt.Sprint("submit happy hour failed, missing city")}
	}

	return nil
}

func (r happyHourSubmitEndpoint) GetPath() string {
	return "/hhsubmit"
}

func (r happyHourSubmitEndpoint) HTTPRequest() interface{} {
	return happyHourSubmitRequest{}
}

func (r happyHourSubmitEndpoint) HTTPResult() interface{} {
	return happyHourSubmitResponse{}
}
