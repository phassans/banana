package controller

import (
	"context"
	"fmt"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
)

type (
	upvoteCountRequest struct {
		ListingID int `json:"listingId"`
	}

	upvoteCountResponse struct {
		upvoteCountRequest
		Count int
		Error *APIError `json:"error,omitempty"`
	}

	upvoteCountEndpoint struct{}
)

var upvoteCount postEndpoint = upvoteCountEndpoint{}

func (r upvoteCountEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(upvoteCountRequest)

	logger := shared.GetLogger()
	logger = logger.With().
		Str("endpoint", r.GetPath()).
		Int("listingId", request.ListingID).Logger()
	logger.Info().Msgf("upvote count request")

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	count, err := rtr.engines.GetUpVotes(request.ListingID)
	result := upvoteCountResponse{upvoteCountRequest: request, Count: count, Error: NewAPIError(err)}
	return result, err
}

func (r upvoteCountEndpoint) Validate(request interface{}) error {
	req := request.(upvoteCountRequest)

	if req.ListingID == 0 {
		return helper.ValidationError{Message: fmt.Sprint("upvote count failed, please provide 'listingId'")}
	}

	return nil
}

func (r upvoteCountEndpoint) GetPath() string {
	return "/upvote/count"
}

func (r upvoteCountEndpoint) HTTPRequest() interface{} {
	return upvoteCountRequest{}
}
