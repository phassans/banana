package controller

import (
	"context"

	"github.com/phassans/banana/model"
)

type (
	favoritesViewRequest struct {
		PhoneID string `json:"phoneId"`
		SortBy  string `json:"sortBy,omitempty"`
	}

	favoritesViewResult struct {
		Result []model.Listing
		Error  *APIError `json:"error,omitempty"`
	}

	favoritesViewEndpoint struct{}
)

var favoritesView postEndpoint = favoritesViewEndpoint{}

func (r favoritesViewEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(favoritesViewRequest)
	result, err := rtr.engines.GetAllFavorites(
		request.PhoneID,
	)
	return favoritesViewResult{Result: result, Error: NewAPIError(err)}, err
}

func (r favoritesViewEndpoint) Validate(request interface{}) error {
	return nil
}

func (r favoritesViewEndpoint) GetPath() string {
	return "/favorites"
}

func (r favoritesViewEndpoint) HTTPRequest() interface{} {
	return favoritesViewRequest{}
}
