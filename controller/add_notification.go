package controller

import (
	"context"

	"github.com/phassans/banana/model"
)

type (
	notificationADDRequest struct {
		PhoneId            string   `json:"phoneId"`
		BusinessId         int      `json:"businessId,omitempty"`
		Price              string   `json:"price,omitempty"`
		Keywords           string   `json:"keywords,omitempty"`
		DietaryRestriction []string `json:"dietaryRestriction,omitempty"`
		Latitude           float64  `json:"latitude,omitempty"`
		Longitude          float64  `json:"longitude,omitempty"`
		Location           string   `json:"location,omitempty"`
	}

	notificationADDResult struct {
		notificationADDRequest
		Error *APIError `json:"error,omitempty"`
	}

	addNotificationEndpoint struct{}
)

var addNotification postEndpoint = addNotificationEndpoint{}

func (r addNotificationEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(notificationADDRequest)

	l := model.Notification{
		PhoneId:            request.PhoneId,
		BusinessId:         request.BusinessId,
		Price:              request.Price,
		Keywords:           request.Keywords,
		DietaryRestriction: request.DietaryRestriction,
		Latitude:           request.Latitude,
		Longitude:          request.Longitude,
		Location:           request.Location,
	}

	err := rtr.engines.AddNotification(l)
	result := notificationADDResult{notificationADDRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r addNotificationEndpoint) Validate(request interface{}) error {
	return nil
}

func (r addNotificationEndpoint) GetPath() string {
	return "/notification/add"
}

func (r addNotificationEndpoint) HTTPRequest() interface{} {
	return notificationADDRequest{}
}
