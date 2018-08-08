package controller

import (
	"context"
	"fmt"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
)

type (
	listingAllRequest struct {
		BusinessID int    `json:"businessId"`
		Status     string `json:"status"`
	}

	listingAllResult struct {
		Result []shared.Listing
		Error  *APIError `json:"error,omitempty"`
	}

	listingAllEndpoint struct{}
)

var listingAll postEndpoint = listingAllEndpoint{}

func (r listingAllEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(listingAllRequest)

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	result, err := rtr.engines.GetListingsByBusinessID(request.BusinessID, request.Status)
	return listingAllResult{Result: result, Error: NewAPIError(err)}, err
}

func (r listingAllEndpoint) Validate(request interface{}) error {
	input := request.(listingAllRequest)

	if input.Status != "all" &&
		input.Status != "active" &&
		input.Status != "scheduled" &&
		input.Status != "ended" {
		return helper.ValidationError{Message: fmt.Sprint("listing all failed, invalid 'status'")}
	}

	if input.BusinessID == 0 {
		return helper.ValidationError{Message: fmt.Sprint("listing all failed, missing businessId")}
	}

	return nil
}

func (r listingAllEndpoint) GetPath() string {
	return "/listing/all"
}

func (r listingAllEndpoint) HTTPRequest() interface{} {
	return listingAllRequest{}
}
