package controller

import (
	"context"

	"github.com/phassans/banana/model"
)

type (
	listingADDRequest struct {
		Title              string   `json:"title"`
		OldPrice           float64  `json:"oldPrice"`
		NewPrice           float64  `json:"newPrice"`
		DietaryRestriction []string `json:"dietaryRestriction,omitempty"`
		Description        string   `json:"description"`
		StartDate          string   `json:"startDate"`
		EndDate            string   `json:"endDate"`
		StartTime          string   `json:"startTime"`
		EndTime            string   `json:"endTime"`
		BusinessID         int      `json:"businessID"`
		Recurring          bool     `json:"recurring"`
		RecurringDays      []string `json:"recurringDays,omitempty"`
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

	l := model.Listing{
		request.Title,
		request.OldPrice,
		request.NewPrice,
		request.DietaryRestriction,
		request.Description,
		request.StartDate,
		request.EndDate,
		request.StartTime,
		request.EndTime,
		request.BusinessID,
		request.Recurring,
		request.RecurringDays,
		1,
	}

	err := rtr.engines.AddListing(l)
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
