package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
)

type (
	favoritesViewRequest struct {
		PhoneID   string  `json:"phoneId"`
		SortBy    string  `json:"sortBy,omitempty"`
		Latitude  float64 `json:"latitude,omitempty"`
		Longitude float64 `json:"longitude,omitempty"`
	}

	favoritesViewResult struct {
		Result []shared.SearchListingResult
		Error  *APIError `json:"error,omitempty"`
	}

	favoritesViewEndpoint struct{}
)

var favouriteAll postEndpoint = favoritesViewEndpoint{}

func (r favoritesViewEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(favoritesViewRequest)

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	result, err := rtr.engines.GetAllFavorites(
		request.PhoneID,
		request.SortBy,
		request.Latitude,
		request.Longitude,
	)
	return favoritesViewResult{Result: result, Error: NewAPIError(err)}, err
}

func (r favoritesViewEndpoint) Validate(request interface{}) error {
	req := request.(favoritesViewRequest)
	if strings.TrimSpace(req.PhoneID) == "" {
		return helper.ValidationError{Message: fmt.Sprint("favorite all failed, please provide 'phoneId'")}
	}

	if strings.TrimSpace(req.SortBy) != "" &&
		strings.TrimSpace(req.SortBy) != "dateAdded" &&
		strings.ToLower(req.SortBy) != "distance" &&
		strings.ToLower(req.SortBy) != "timeleft" &&
		strings.ToLower(req.SortBy) != "price" {
		return helper.ValidationError{Message: fmt.Sprint("favorite all failed, invalid 'sortBy'")}
	}

	return nil
}

func (r favoritesViewEndpoint) GetPath() string {
	return "/favorite/all"
}

func (r favoritesViewEndpoint) HTTPRequest() interface{} {
	return favoritesViewRequest{}
}
