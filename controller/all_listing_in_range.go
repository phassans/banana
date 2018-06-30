package controller

import (
	"context"

	"github.com/pshassans/banana/model"
)

type (
	allListingInRangeRequest struct {
		Latitude      float64 `json:"latitude,omitempty"`
		Longitude     float64 `json:"longitude,omitempty"`
		ZipCode       int     `json:"zipCode,omitempty"`
		PriceFilter   string  `json:"priceFilter,omitempty"`
		DietaryFilter string  `json:"dietaryFilter,omitempty"`
		Keywords      string  `json:"keywords,omitempty"`
	}

	allListingInRangeResult struct {
		Result []model.Listing
		Error  *APIError `json:"error,omitempty"`
	}

	allListingInRangeEndpoint struct{}
)

var allListing postEndpoint = allListingInRangeEndpoint{}

func (r allListingInRangeEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(allListingInRangeRequest)
	result, err := rtr.engines.GetAllListingsInRange(
		request.Latitude,
		request.Longitude,
		request.ZipCode,
		request.PriceFilter,
		request.DietaryFilter,
		request.Keywords,
	)
	return allListingInRangeResult{Result: result, Error: NewAPIError(err)}, err
}

func (r allListingInRangeEndpoint) Validate(request interface{}) error {
	return nil
}

func (r allListingInRangeEndpoint) GetPath() string {
	return "/listing/all"
}

func (r allListingInRangeEndpoint) HTTPRequest() interface{} {
	return allListingInRangeRequest{}
}
