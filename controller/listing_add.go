package controller

import (
	"context"

	"github.com/phassans/banana/shared"
)

type (
	listingADDRequest struct {
		BusinessID         int      `json:"businessID"`
		Title              string   `json:"title"`
		OldPrice           float64  `json:"oldPrice,omitempty"`
		NewPrice           float64  `json:"newPrice"`
		Discount           float64  `json:"discount,omitempty"`
		DietaryRestriction []string `json:"dietaryRestriction,omitempty"`
		Description        string   `json:"description"`
		StartDate          string   `json:"startDate"`
		StartTime          string   `json:"startTime"`
		EndTime            string   `json:"endTime,omitempty"`
		MultipleDays       bool     `json:"multipleDays"`
		EndDate            string   `json:"endDate,omitempty"`
		Recurring          bool     `json:"recurring"`
		RecurringDays      []string `json:"recurringDays,omitempty"`
		RecurringEndDate   string   `json:"recurringEndDate,omitempty"`
		Type               string   `json:"type"`
	}

	listingADDResult struct {
		listingADDRequest
		Error *APIError `json:"error,omitempty"`
	}

	addListingEndpoint struct{}
)

var addListing postEndpoint = addListingEndpoint{}

func (r addListingEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(listingADDRequest)

	l := shared.Listing{
		Title:              request.Title,
		OldPrice:           request.OldPrice,
		NewPrice:           request.NewPrice,
		Discount:           request.Discount,
		DietaryRestriction: request.DietaryRestriction,
		Description:        request.Description,
		StartDate:          request.StartDate,
		StartTime:          request.StartTime,
		EndTime:            request.EndTime,
		BusinessID:         request.BusinessID,
		MultipleDays:       request.MultipleDays,
		EndDate:            request.EndDate,
		Recurring:          request.Recurring,
		RecurringDays:      request.RecurringDays,
		RecurringEndDate:   request.RecurringEndDate,
		Type:               request.Type,
	}

	err := rtr.engines.AddListing(&l)
	result := listingADDResult{listingADDRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r addListingEndpoint) Validate(request interface{}) error {
	return nil
}

func (r addListingEndpoint) GetPath() string {
	return "/listing/add"
}

func (r addListingEndpoint) HTTPRequest() interface{} {
	return listingADDRequest{}
}
