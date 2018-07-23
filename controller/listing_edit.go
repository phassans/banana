package controller

import (
	"context"

	"github.com/phassans/banana/shared"
	"github.com/rs/xlog"
)

type (
	listingEditRequest struct {
		ListingID          int      `json:"listingId"`
		BusinessID         int      `json:"businessId"`
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
		Type               string   `json:"listingType"`
	}

	listingEditResult struct {
		listingEditRequest
		Error *APIError `json:"error,omitempty"`
	}

	editListingEndpoint struct{}
)

var listingEdit postEndpoint = editListingEndpoint{}

func (r editListingEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(listingEditRequest)
	xlog.Infof("POST %s query %+v", r.GetPath(), request)

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
		ListingID:          request.ListingID,
	}

	err := rtr.engines.ListingEdit(&l)
	result := listingEditResult{listingEditRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r editListingEndpoint) Validate(request interface{}) error {
	return nil
}

func (r editListingEndpoint) GetPath() string {
	return "/listing/edit"
}

func (r editListingEndpoint) HTTPRequest() interface{} {
	return listingEditRequest{}
}
