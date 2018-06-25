package api

import (
	"context"
)

type businessRequest struct {
	BusinessName string `json:"name"`
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
	_, err := rtr.engine.CreateBusiness(request.BusinessName)
	result := businessResult{businessRequest: request, Error: NewAPIError(err)}
	return result, nil
}

var createBusiness postEndpoint = createBusinessEndpoint{}
