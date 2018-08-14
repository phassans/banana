package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/banana/helper"
)

type (
	notificationADDRequest struct {
		PhoneID        string   `json:"phoneId"`
		Latitude       float64  `json:"latitude,omitempty"`
		Longitude      float64  `json:"longitude,omitempty"`
		Location       string   `json:"location,omitempty"`
		PriceFilter    string   `json:"priceFilter,omitempty"`
		DietaryFilters []string `json:"dietaryFilters,omitempty"`
		DistanceFilter string   `json:"distanceFilter,omitempty"`
		Keywords       string   `json:"keywords,omitempty"`
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

	// validate input
	if err := r.Validate(request); err != nil {
		return nil, err
	}

	err := rtr.engines.AddNotification(
		request.PhoneID,
		request.Latitude,
		request.Longitude,
		request.Location,
		request.PriceFilter,
		request.DietaryFilters,
		request.DistanceFilter,
		request.Keywords,
	)
	result := notificationADDResult{notificationADDRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r addNotificationEndpoint) Validate(request interface{}) error {
	req := request.(notificationADDRequest)

	if strings.TrimSpace(req.PhoneID) == "" {
		return helper.ValidationError{Message: fmt.Sprint("notification add failed, please provide 'phoneId'")}
	}

	if strings.TrimSpace(req.Location) == "" && (req.Latitude == 0 || req.Longitude == 0) {
		return helper.ValidationError{Message: fmt.Sprint("notification add failed, please provide 'location' or 'latitude' and 'longitude'")}
	}

	for _, dietary := range req.DietaryFilters {
		if strings.ToLower(dietary) != "vegetarian" &&
			strings.ToLower(dietary) != "vegan" &&
			strings.ToLower(dietary) != "gluten free" {
			return helper.ValidationError{Message: fmt.Sprint("notification add failed, invalid dietaryRestriction")}
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
