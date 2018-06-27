package api

import (
	"context"
)

type businessRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
}

type businessResult struct {
	businessRequest
	Error *APIError `json:"error,omitempty"`
}

type createBusinessEndpoint struct{}

func (r createBusinessEndpoint) GetPath() string          { return "/business/create" }
func (r createBusinessEndpoint) HTTPRequest() interface{} { return businessRequest{} }
func (r createBusinessEndpoint) HTTPResult() interface{}  { return businessResult{} }
func (r createBusinessEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(businessRequest)
	_, err := rtr.engines.AddBusiness(request.Name, request.Phone, request.Website)
	result := businessResult{businessRequest: request, Error: NewAPIError(err)}
	return result, nil
}

var addBusiness postEndpoint = createBusinessEndpoint{}
