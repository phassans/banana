package controller

import (
	"context"

	"github.com/pshassans/banana/model"
)

type allListingInRangeRequest struct {
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	ZipCode       int     `json:"zipCode,omitempty"`
	PriceFilter   string  `json:"priceFilter,omitempty"`
	DietaryFilter string  `json:"dietaryFilter,omitempty"`
	Keywords      string  `json:"keywords,omitempty"`
}

type allListingInRangeResult struct {
	Result []model.Listing
	Error  *APIError `json:"error,omitempty"`
}

type allListingInRangeEndpoint struct{}

func (r allListingInRangeEndpoint) GetPath() string          { return "/listing/all" }
func (r allListingInRangeEndpoint) HTTPRequest() interface{} { return allListingInRangeRequest{} }
func (r allListingInRangeEndpoint) HTTPResult() interface{}  { return allListingInRangeResult{} }
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

var allListing postEndpoint = allListingInRangeEndpoint{}
