package controller

import (
	"context"

	"fmt"
	"strings"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
	"github.com/rs/xlog"
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
	xlog.Infof("POST %s query %+v", r.GetPath(), request)

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	result, err := rtr.engines.GetAllNotifications(
		request.PhoneID,
	)
	return notificationAllResult{Result: result, Error: NewAPIError(err)}, err
}

func (r notificationAllEndpoint) Validate(request interface{}) error {
	req := request.(notificationAllRequest)
	if strings.TrimSpace(req.PhoneID) == "" {
		return helper.ValidationError{Message: fmt.Sprint("notification all failed, please provide 'phoneId'")}
	}

	return nil
}

func (r notificationAllEndpoint) GetPath() string {
	return "/notification/all"
}

func (r notificationAllEndpoint) HTTPRequest() interface{} {
	return notificationAllRequest{}
}
