package controller

import (
	"context"
	"fmt"

	"github.com/phassans/banana/helper"
)

type (
	deleteBusinessRequest struct {
		BusinessID int `json:"businessId"`
	}

	deleteBusinessResponse struct {
		deleteBusinessRequest
		Error *APIError `json:"error,omitempty"`
	}

	deleteBusinessEndpoint struct{}
)

var businessDelete postEndpoint = deleteBusinessEndpoint{}

func (r deleteBusinessEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(deleteBusinessRequest)

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	err := rtr.engines.BusinessDelete(request.BusinessID)
	result := deleteBusinessResponse{deleteBusinessRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r deleteBusinessEndpoint) Validate(request interface{}) error {
	input := request.(deleteBusinessRequest)
	if input.BusinessID == 0 {
		return helper.ValidationError{Message: fmt.Sprint("business delete failed, missing businessId")}
	}

	return nil
}

func (r deleteBusinessEndpoint) GetPath() string {
	return "/business/delete"
}

func (r deleteBusinessEndpoint) HTTPRequest() interface{} {
	return deleteBusinessRequest{}
}
