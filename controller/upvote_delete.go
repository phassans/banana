package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
)

type (
	upvoteDeleteRequest struct {
		PhoneID   string `json:"phoneId"`
		ListingID int    `json:"listingId"`
	}

	upvoteDeleteResponse struct {
		upvoteDeleteRequest
		Error *APIError `json:"error,omitempty"`
	}

	upvoteDeleteEndpoint struct{}
)

var upvoteDelete postEndpoint = upvoteDeleteEndpoint{}

func (r upvoteDeleteEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(upvoteDeleteRequest)

	logger := shared.GetLogger()
	logger = logger.With().
		Str("endpoint", r.GetPath()).
		Str("phoneId", request.PhoneID).
		Int("listingId", request.ListingID).Logger()
	logger.Info().Msgf("upvote delete request")

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	err := rtr.engines.DeleteUpVote(request.PhoneID, request.ListingID)
	result := upvoteDeleteResponse{upvoteDeleteRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r upvoteDeleteEndpoint) Validate(request interface{}) error {
	req := request.(upvoteDeleteRequest)

	if strings.TrimSpace(req.PhoneID) == "" {
		return helper.ValidationError{Message: fmt.Sprint("upvote delete failed, please provide 'phoneId'")}
	}

	if req.ListingID == 0 {
		return helper.ValidationError{Message: fmt.Sprint("upvote delete failed, please provide 'listingId'")}
	}

	return nil
}

func (r upvoteDeleteEndpoint) GetPath() string {
	return "/upvote/delete"
}

func (r upvoteDeleteEndpoint) HTTPRequest() interface{} {
	return upvoteDeleteRequest{}
}
