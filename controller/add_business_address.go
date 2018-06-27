package controller

import (
	"context"
)

type businessAddressRequest struct {
	Line1        string `json:"line1"`
	Line2        string `json:"line2"`
	City         string `json:"city"`
	PostalCode   string `json:"postalCode"`
	State        string `json:"state"`
	Country      string `json:"country"`
	BusinessName string `json:"businessName"`
	OtherDetails string `json:"otherDetails"`
}

type businessAddressResult struct {
	businessAddressRequest
	Error *APIError `json:"error,omitempty"`
}

type addBusinessAddressEndpoint struct{}

func (r addBusinessAddressEndpoint) GetPath() string          { return "/business/address/add" }
func (r addBusinessAddressEndpoint) HTTPRequest() interface{} { return businessAddressRequest{} }
func (r addBusinessAddressEndpoint) HTTPResult() interface{}  { return businessAddressResult{} }
func (r addBusinessAddressEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(businessAddressRequest)
	err := rtr.engines.AddBusinessAddress(
		request.Line1,
		request.Line2,
		request.City,
		request.PostalCode,
		request.State,
		request.Country,
		request.BusinessName,
		request.OtherDetails,
	)
	result := businessAddressResult{businessAddressRequest: request, Error: NewAPIError(err)}
	return result, nil
}

var addBusinessAddress postEndpoint = addBusinessAddressEndpoint{}
