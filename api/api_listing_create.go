package api

import (
	"context"

	"github.com/banana/helper"
	"github.com/rs/xlog"
)

type listingRequest struct {
	ListingName string `json:"listingName"`
}

type listingResult struct {
	listingRequest
	Error *APIError `json:"error,omitempty"`
}

type createListingEndpoint struct{}

func (r createListingEndpoint) GetPath() string          { return "/listing/create" }
func (r createListingEndpoint) HTTPRequest() interface{} { return listingRequest{} }
func (r createListingEndpoint) HTTPResult() interface{}  { return listingResult{} }
func (r createListingEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(listingRequest)
	err := CreateListing(ctx, request.ListingName)
	result := listingResult{listingRequest: request, Error: NewAPIError(err)}
	return result, nil
}

func CreateListing(ctx context.Context, name string) error {
	xlog.Infof("whoohoo logger %s with apiVersion %s", name, helper.GetContextValue(ctx, helper.ApiVersion))
	return nil
}

var createListing postEndpoint = createListingEndpoint{}
