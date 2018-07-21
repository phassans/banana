package controller

import (
	"context"

	"github.com/phassans/banana/shared"
)

type (
	notificationAllRequest struct {
		PhoneID string `json:"phoneId"`
	}

	notificationAllResult struct {
		Result []shared.Notification
		Error  *APIError `json:"error,omitempty"`
	}

	notificationAllEndpoint struct{}
)

var notificationAll postEndpoint = notificationAllEndpoint{}

func (r notificationAllEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(notificationAllRequest)
	result, err := rtr.engines.GetAllNotifications(
		request.PhoneID,
	)
	return notificationAllResult{Result: result, Error: NewAPIError(err)}, err
}

func (r notificationAllEndpoint) Validate(request interface{}) error {
	return nil
}

func (r notificationAllEndpoint) GetPath() string {
	return "/notification/all"
}

func (r notificationAllEndpoint) HTTPRequest() interface{} {
	return notificationAllRequest{}
}
