package controller

import (
	"context"
)

type (
	businessRequest struct {
		Name       string `json:"name"`
		Phone      string `json:"phone"`
		Website    string `json:"website"`
		Street     string `json:"street"`
		City       string `json:"city"`
		PostalCode string `json:"postalCode"`
		State      string `json:"state"`
		Country    string `json:"country"`
	}

	businessResult struct {
		businessRequest
		Error *APIError `json:"error,omitempty"`
	}

	createBusinessEndpoint struct{}
)

var addBusiness postEndpoint = createBusinessEndpoint{}

func (r createBusinessEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(businessRequest)
	_, err := rtr.engines.AddBusiness(
		request.Name,
		request.Phone,
		request.Website,
		request.Street,
		request.City,
		request.PostalCode,
		request.State,
		request.Country,
	)
	result := businessResult{businessRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r createBusinessEndpoint) Validate(request interface{}) error {
	return nil
}

func (r createBusinessEndpoint) GetPath() string {
	return "/business/add"
}

func (r createBusinessEndpoint) HTTPRequest() interface{} {
	return businessRequest{}
}
