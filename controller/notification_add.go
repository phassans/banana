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
	notificationADDRequest struct {
		PhoneId            string   `json:"phoneId"`
		BusinessId         int      `json:"businessId"`
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

var notificationAdd postEndpoint = addNotificationEndpoint{}

func (r addNotificationEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(notificationADDRequest)
	xlog.Infof("POST %s query %+v", r.GetPath(), request)

	// validate input
	if err := r.Validate(request); err != nil {
		return nil, err
	}

	l := shared.Notification{
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
	req := request.(notificationADDRequest)

	if strings.TrimSpace(req.PhoneId) == "" {
		return helper.ValidationError{Message: fmt.Sprint("notification add failed, please provide 'phoneId'")}
	}

	if req.BusinessId == 0 {
		return helper.ValidationError{Message: fmt.Sprint("notification add failed, please provide 'businessId'")}
	}

	if strings.TrimSpace(req.Location) == "" && (req.Latitude == 0 || req.Longitude == 0) {
		return helper.ValidationError{Message: fmt.Sprint("listings search failed, please provide 'location' or 'latitude' and 'longitude'")}
	}

	for _, dietary := range req.DietaryRestriction {
		if strings.ToLower(dietary) != "vegetarian" &&
			strings.ToLower(dietary) != "vegan" &&
			strings.ToLower(dietary) != "gluten free" {
			return helper.ValidationError{Message: fmt.Sprint("listings search failed, invalid dietaryRestriction")}
		}
	}

	return nil
}

func (r addNotificationEndpoint) GetPath() string {
	return "/notification/add"
}

func (r addNotificationEndpoint) HTTPRequest() interface{} {
	return notificationADDRequest{}
}
