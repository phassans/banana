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
		ListingType        string   `json:"listingType"`
		ImageLink          string   `json:"imageLink,omitempty"`
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

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

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
		Type:               request.ListingType,
		ListingID:          request.ListingID,
		ImageLink:          request.ImageLink,
	}

	err := rtr.engines.ListingEdit(&l)
	result := listingEditResult{listingEditRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r editListingEndpoint) Validate(request interface{}) error {
	input := request.(listingEditRequest)

	if input.ListingType != shared.ListingTypeMeal && input.ListingType != shared.ListingTypeHappyHour {
		return helper.ValidationError{Message: fmt.Sprint("listing edit failed, invalid 'listingType'")}
	}

	var businessFields = []string{input.Title, input.Description, input.StartDate, input.StartTime, input.EndTime}
	for _, field := range businessFields {
		if strings.TrimSpace(field) == "" {
			return helper.ValidationError{Message: fmt.Sprint("listing edit failed, missing mandatory fields")}
		}
	}

	if input.ListingType == shared.ListingTypeMeal && input.NewPrice == 0 {
		return helper.ValidationError{Message: fmt.Sprint("listing edit failed, add 'newPrice' for the meal")}
	}

	if input.ListingType == shared.ListingTypeHappyHour && (input.Discount == 0 && input.NewPrice == 0) {
		return helper.ValidationError{Message: fmt.Sprint("listing edit failed, add 'discount' for the happyhour")}
	}

	if input.MultipleDays && input.Recurring {
		return helper.ValidationError{Message: fmt.Sprint("listing edit failed, listing cannot be multiple days and recurring")}
	}

	if input.MultipleDays && input.EndDate == "" {
		return helper.ValidationError{Message: fmt.Sprint("listing edit failed, please provide 'endDate' for multiple days lising")}
	}

	if input.Recurring && input.RecurringEndDate == "" {
		return helper.ValidationError{Message: fmt.Sprint("listing edit failed, please provide 'recurringEndDate' for recurring listing")}
	} else if input.Recurring && len(input.RecurringDays) == 0 {
		return helper.ValidationError{Message: fmt.Sprint("listing edit failed, please provide 'recurringDays' for recurring listing")}
	}

	if input.ListingID == 0 {
		return helper.ValidationError{Message: fmt.Sprint("listing edit failed, please provide 'listingId' for listing")}
	}

	if input.BusinessID == 0 {
		return helper.ValidationError{Message: fmt.Sprint("listing edit failed, please provide 'businessId' for listing")}
	}

	return nil
}

func (r editListingEndpoint) GetPath() string {
	return "/listing/edit"
}

func (r editListingEndpoint) HTTPRequest() interface{} {
	return listingEditRequest{}
}
