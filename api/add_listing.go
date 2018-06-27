package api

import (
	"context"
)

type listingADDRequest struct {
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	OldPrice     float64 `json:"oldPrice"`
	NewPrice     float64 `json:"newPrice"`
	ListingDate  string  `json:"listingDate"`
	StartTime    string  `json:"startTime"`
	EndTime      string  `json:"endTime"`
	BusinessName string  `json:"businessName"`
	Recurring    bool    `json:"recurring"`
}

type listingADDResult struct {
	listingADDRequest
	Error *APIError `json:"error,omitempty"`
}

type addListingEndpoint struct{}

func (r addListingEndpoint) GetPath() string          { return "/listing/add" }
func (r addListingEndpoint) HTTPRequest() interface{} { return listingADDRequest{} }
func (r addListingEndpoint) HTTPResult() interface{}  { return listingADDResult{} }
func (r addListingEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(listingADDRequest)
	err := rtr.engines.AddListing(
		request.Title,
		request.Description,
		request.OldPrice,
		request.NewPrice,
		request.ListingDate,
		request.StartTime,
		request.EndTime,
		request.BusinessName,
	)
	result := listingADDResult{listingADDRequest: request, Error: NewAPIError(err)}
	return result, nil
}

var addListing postEndpoint = addListingEndpoint{}
