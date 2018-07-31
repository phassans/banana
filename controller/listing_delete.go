package controller

import (
	"context"
	"fmt"

	"github.com/phassans/banana/helper"
	"github.com/rs/xlog"
)

type (
	deleteListingRequest struct {
		ListingID int `json:"listingId"`
	}

	deleteListingResponse struct {
		deleteListingRequest
		Error *APIError `json:"error,omitempty"`
	}

	deleteListingEndpoint struct{}
)

var listingDelete postEndpoint = deleteListingEndpoint{}

func (r deleteListingEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(deleteListingRequest)
	xlog.Infof("POST %s query %+v", r.GetPath(), request)

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	err := rtr.engines.DeleteListing(request.ListingID)
	result := deleteListingResponse{deleteListingRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r deleteListingEndpoint) Validate(request interface{}) error {
	input := request.(deleteListingRequest)
	if input.ListingID == 0 {
		return helper.ValidationError{Message: fmt.Sprint("listing all failed, missing listingId")}
	}

	return nil
}

func (r deleteListingEndpoint) GetPath() string {
	return "/listing/delete"
}

func (r deleteListingEndpoint) HTTPRequest() interface{} {
	return deleteListingRequest{}
}
