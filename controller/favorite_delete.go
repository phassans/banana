package controller

import (
	"context"
)

type (
	deleteFavoriteRequest struct {
		PhoneID   string `json:"phoneId"`
		ListingID int    `json:"listingId"`
	}

	deleteFavoriteResponse struct {
		deleteFavoriteRequest
		Error *APIError `json:"error,omitempty"`
	}

	deleteFavoriteEndpoint struct{}
)

var deleteFavorite postEndpoint = deleteFavoriteEndpoint{}

func (r deleteFavoriteEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(deleteFavoriteRequest)
	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	err := rtr.engines.DeleteFavorite(request.PhoneID, request.ListingID)
	result := deleteFavoriteResponse{deleteFavoriteRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r deleteFavoriteEndpoint) Validate(request interface{}) error {
	return nil
}

func (r deleteFavoriteEndpoint) GetPath() string {
	return "/favorite/delete"
}

func (r deleteFavoriteEndpoint) HTTPRequest() interface{} {
	return deleteFavoriteRequest{}
}
