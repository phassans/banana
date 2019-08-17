package controller

import (
	"context"
)

type (
	webhookRequest struct {
		// Business
		Name       string `json:"name"`
		Phone      string `json:"phone"`
		Website    string `json:"website,omitempty"`
		Street     string `json:"street"`
		City       string `json:"city"`
		PostalCode string `json:"postalCode"`
		State      string `json:"state"`
		Monday     string `json:"monday"`
		Tuesday    string `json:"tuesday"`

		// Listing
		Title         string `json:"title"`
		Description   string `json:"description"`
		StartDate     string `json:"startDate"`
		EndDate       string `json:"endDate"`
		RecurringDays string `json:"recurringDays"`
		StartTime     string `json:"startTime"`
		EndTime       string `json:"endTime"`
		AddTime       string `json:"addTime"`
	}

	webhookResult struct {
		webhookRequest
		Error *APIError `json:"error,omitempty"`
	}

	webhookEndpoint struct{}
)

var webhook postEndpoint = webhookEndpoint{}

func (r webhookEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(webhookRequest)

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	result := webhookResult{webhookRequest: request}
	return result, nil
}

func (r webhookEndpoint) Validate(request interface{}) error {
	return nil
}

func (r webhookEndpoint) GetPath() string {
	return "/webhook"
}

func (r webhookEndpoint) HTTPRequest() interface{} {
	return webhookRequest{}
}
