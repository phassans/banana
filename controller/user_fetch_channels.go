package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/banana/shared"

	"github.com/phassans/banana/helper"
)

type (
	fetchUserChannelsRequest struct {
		UserID    string `json:"userId"`
		Education bool   `json:"education,omitempty"`
		Company   bool   `json:"company,omitempty"`
	}

	fetchUserChannelsResponse struct {
		fetchUserChannelsRequest
		Error    *APIError       `json:"error,omitempty"`
		Channels shared.Channels `json:"channels,omitempty"`
	}

	fetchUserChannelsEndpoint struct{}
)

var fetchUserChannels postEndpoint = fetchUserChannelsEndpoint{}

func (r fetchUserChannelsEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(fetchUserChannelsRequest)
	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	chans, _ := rtr.engines.FetchUserChannels(request.UserID)
	return fetchUserChannelsResponse{fetchUserChannelsRequest: request, Error: NewAPIError(nil), Channels: chans}, nil

	result := fetchUserChannelsResponse{fetchUserChannelsRequest: request, Error: NewAPIError(nil)}
	return result, nil
}

func (r fetchUserChannelsEndpoint) Validate(request interface{}) error {
	input := request.(fetchUserChannelsRequest)
	if strings.TrimSpace(input.UserID) == "" {
		return helper.ValidationError{Message: fmt.Sprint("fetch user education failed, missing fields")}
	}

	return nil
}

func (r fetchUserChannelsEndpoint) GetPath() string {
	return "/user/channels"
}

func (r fetchUserChannelsEndpoint) HTTPRequest() interface{} {
	return fetchUserChannelsRequest{}
}

func (r fetchUserChannelsEndpoint) HTTPResult() interface{} {
	return fetchUserChannelsResponse{}
}
