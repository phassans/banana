package controller

import (
	"context"

	"github.com/phassans/banana/model"
)

type (
	listingsSearchRequest struct {
		ListingType   string  `json:"listingType"`
		Future        bool    `json:"future"`
		Latitude      float64 `json:"latitude,omitempty"`
		Longitude     float64 `json:"longitude,omitempty"`
		ZipCode       int     `json:"zipCode,omitempty"`
		PriceFilter   string  `json:"priceFilter,omitempty"`
		DietaryFilter string  `json:"dietaryFilter,omitempty"`
		Keywords      string  `json:"keywords,omitempty"`
		SortBy        string  `json:"sortBy,omitempty"`
	}

	listingsSearchResult struct {
		Result []model.Listing
		Error  *APIError `json:"error,omitempty"`
	}

	listingsSearchEndpoint struct{}
)

var listingsSearch postEndpoint = listingsSearchEndpoint{}

func (r listingsSearchEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(listingsSearchRequest)
	result, err := rtr.engines.SearchListings(
		request.ListingType,
		request.Latitude,
		request.Longitude,
		request.ZipCode,
		request.PriceFilter,
		request.DietaryFilter,
		request.Keywords,
		request.SortBy,
	)
	return listingsSearchResult{Result: result, Error: NewAPIError(err)}, err
}

func (r listingsSearchEndpoint) Validate(request interface{}) error {
	return nil
}

func (r listingsSearchEndpoint) GetPath() string {
	return "/listings/search"
}

func (r listingsSearchEndpoint) HTTPRequest() interface{} {
	return listingsSearchRequest{}
}
