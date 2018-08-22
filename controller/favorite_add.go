package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
)

type (
	addFavoriteRequest struct {
		PhoneID       string `json:"phoneId"`
		ListingID     int    `json:"listingId"`
		ListingDateID int    `json:"listingDateId"`
	}

	addFavoriteResponse struct {
		addFavoriteRequest
		Error *APIError `json:"error,omitempty"`
	}

	addFavoriteEndpoint struct{}
)

var favouriteAdd postEndpoint = addFavoriteEndpoint{}

func (r addFavoriteEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(addFavoriteRequest)

	logger := shared.GetLogger()
	logger = logger.With().
		Str("endpoint", r.GetPath()).
		Str("phoneId", request.PhoneID).
		Int("listingId", request.ListingID).
		Int("listingDateId", request.ListingDateID).Logger()
	logger.Info().Msgf("favorite add request")

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	err := rtr.engines.AddFavorite(request.PhoneID, request.ListingID, request.ListingDateID)
	result := addFavoriteResponse{addFavoriteRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r addFavoriteEndpoint) Validate(request interface{}) error {
	req := request.(addFavoriteRequest)

	if strings.TrimSpace(req.PhoneID) == "" {
		return helper.ValidationError{Message: fmt.Sprint("favorite add failed, please provide 'phoneId'")}
	}

	if req.ListingID == 0 {
		return helper.ValidationError{Message: fmt.Sprint("favorite add failed, please provide 'listingId'")}
	}

	return nil
}

func (r addFavoriteEndpoint) GetPath() string {
	return "/favorite/add"
}

func (r addFavoriteEndpoint) HTTPRequest() interface{} {
	return addFavoriteRequest{}
}
