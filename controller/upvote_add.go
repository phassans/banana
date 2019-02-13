package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
)

type (
	upvoteAddRequest struct {
		PhoneID   string `json:"phoneId"`
		ListingID int    `json:"listingId"`
	}

	upvoteAddResponse struct {
		upvoteAddRequest
		Error *APIError `json:"error,omitempty"`
	}

	upvoteAddEndpoint struct{}
)

var upvoteAdd postEndpoint = upvoteAddEndpoint{}

func (r upvoteAddEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(upvoteAddRequest)

	logger := shared.GetLogger()
	logger = logger.With().
		Str("endpoint", r.GetPath()).
		Str("phoneId", request.PhoneID).
		Int("listingId", request.ListingID).Logger()
	logger.Info().Msgf("upvote add request")

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	err := rtr.engines.AddUpVote(request.PhoneID, request.ListingID)
	result := upvoteAddResponse{upvoteAddRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r upvoteAddEndpoint) Validate(request interface{}) error {
	req := request.(upvoteAddRequest)

	if strings.TrimSpace(req.PhoneID) == "" {
		return helper.ValidationError{Message: fmt.Sprint("upvote add failed, please provide 'phoneId'")}
	}

	if req.ListingID == 0 {
		return helper.ValidationError{Message: fmt.Sprint("upvote add failed, please provide 'listingId'")}
	}

	return nil
}

func (r upvoteAddEndpoint) GetPath() string {
	return "/upvote/add"
}

func (r upvoteAddEndpoint) HTTPRequest() interface{} {
	return upvoteAddRequest{}
}
