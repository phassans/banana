package api

import (
	"context"
)

type ownerRequest struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	OwnerPhone   string `json:"phone"`
	Email        string `json:"email"`
	BusinessName string `json:"businessName"`
}

type ownerResult struct {
	ownerRequest
	Error *APIError `json:"error,omitempty"`
}

type addOwnerEndpoint struct{}

func (r addOwnerEndpoint) GetPath() string          { return "/owner/add" }
func (r addOwnerEndpoint) HTTPRequest() interface{} { return ownerRequest{} }
func (r addOwnerEndpoint) HTTPResult() interface{}  { return ownerResult{} }
func (r addOwnerEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(ownerRequest)
	err := rtr.ownerEngine.AddOwner(request.FirstName, request.LastName, request.OwnerPhone, request.Email, request.BusinessName)
	result := ownerResult{ownerRequest: request, Error: NewAPIError(err)}
	return result, nil
}

var addOwner postEndpoint = addOwnerEndpoint{}
