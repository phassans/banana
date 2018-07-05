package controller

import (
	"context"
)

type (
	addFavoriteRequest struct {
		PhoneID   string `json:"phoneId"`
		ListingID int    `json:"listingId"`
	}

	addFavoriteResponse struct {
		addFavoriteRequest
		Error *APIError `json:"error,omitempty"`
	}

	addFavoriteEndpoint struct{}
)

var addFavorite postEndpoint = addFavoriteEndpoint{}

func (r addFavoriteEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(addFavoriteRequest)
	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	err := rtr.engines.AddFavorite(request.PhoneID, request.ListingID)
	result := addFavoriteResponse{addFavoriteRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r addFavoriteEndpoint) Validate(request interface{}) error {
	return nil
}

func (r addFavoriteEndpoint) GetPath() string {
	return "/favorite/add"
}

func (r addFavoriteEndpoint) HTTPRequest() interface{} {
	return addFavoriteRequest{}
}
