package controller

import (
	"context"
	"fmt"

	"github.com/phassans/banana/helper"
)

type (
	deleteNotificationRequest struct {
		NotificationID int `json:"notificationId"`
	}

	deleteNotificationResponse struct {
		deleteNotificationRequest
		Error *APIError `json:"error,omitempty"`
	}

	deleteNotificationEndpoint struct{}
)

var notificationDelete postEndpoint = deleteNotificationEndpoint{}

func (r deleteNotificationEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(deleteNotificationRequest)
	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	err := rtr.engines.DeleteNotification(request.NotificationID)
	result := deleteNotificationResponse{deleteNotificationRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r deleteNotificationEndpoint) Validate(request interface{}) error {
	req := request.(deleteNotificationRequest)
	if req.NotificationID == 0 {
		return helper.ValidationError{Message: fmt.Sprint("notification delete failed, please provide 'notificationId'")}
	}
	return nil
}

func (r deleteNotificationEndpoint) GetPath() string {
	return "/notification/delete"
}

func (r deleteNotificationEndpoint) HTTPRequest() interface{} {
	return deleteNotificationRequest{}
}
