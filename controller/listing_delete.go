package controller

import (
	"context"
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
	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	err := rtr.engines.DeleteListing(request.ListingID)
	result := deleteListingResponse{deleteListingRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r deleteListingEndpoint) Validate(request interface{}) error {
	return nil
}

func (r deleteListingEndpoint) GetPath() string {
	return "/listing/delete"
}

func (r deleteListingEndpoint) HTTPRequest() interface{} {
	return deleteListingRequest{}
}