package controller

import (
	"context"

	"github.com/phassans/banana/shared"
)

type (
	listingAllRequest struct {
		BusinessID int    `json:"businessID"`
		status     string `json:"status"`
	}

	listingAllResult struct {
		Result []shared.Listing
		Error  *APIError `json:"error,omitempty"`
	}

	listingAllEndpoint struct{}
)

var listingAll postEndpoint = listingAllEndpoint{}

func (r listingAllEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(listingAllRequest)
	result, err := rtr.engines.GetAllListings(request.BusinessID, request.status)
	return listingAllResult{Result: result, Error: NewAPIError(err)}, err
}

func (r listingAllEndpoint) Validate(request interface{}) error {
	return nil
}

func (r listingAllEndpoint) GetPath() string {
	return "/listing/all"
}

func (r listingAllEndpoint) HTTPRequest() interface{} {
	return listingAllRequest{}
}
