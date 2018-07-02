package controller

import (
	"context"
)

type (
	listingAllRequest struct {
		BusinessID int `json:"businessID"`
	}

	allListingEndpoint struct{}
)

var allListingAdmin postEndpoint = allListingEndpoint{}

func (r allListingEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(listingAllRequest)

	return rtr.engines.GetAllListings(request.BusinessID)
}

func (r allListingEndpoint) Validate(request interface{}) error {
	return nil
}

func (r allListingEndpoint) GetPath() string {
	return "/listing/all"
}

func (r allListingEndpoint) HTTPRequest() interface{} {
	return listingAllRequest{}
}
