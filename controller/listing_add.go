package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
)

type (
	listingADDRequest struct {
		BusinessID          int      `json:"businessId"`
		Title               string   `json:"title"`
		OldPrice            float64  `json:"oldPrice,omitempty"`
		NewPrice            float64  `json:"newPrice"`
		Discount            float64  `json:"discount,omitempty"`
		DiscountDescription string   `json:"discountDescription,omitempty"`
		DietaryRestriction  []string `json:"dietaryRestriction,omitempty"`
		Description         string   `json:"description"`
		StartDate           string   `json:"startDate"`
		StartTime           string   `json:"startTime"`
		EndTime             string   `json:"endTime,omitempty"`
		MultipleDays        bool     `json:"multipleDays"`
		EndDate             string   `json:"endDate,omitempty"`
		Recurring           bool     `json:"recurring"`
		RecurringDays       []string `json:"recurringDays,omitempty"`
		RecurringEndDate    string   `json:"recurringEndDate,omitempty"`
		ListingType         string   `json:"listingType"`
		ImageLink           string   `json:"imageLink,omitempty"`

		ListingID int `json:"listingId,omitempty"`
	}

	listingADDResult struct {
		listingADDRequest
		Error *APIError `json:"error,omitempty"`
	}

	addListingEndpoint struct{}
)

var listingAdd postEndpoint = addListingEndpoint{}

func (r addListingEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(listingADDRequest)

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	l := shared.Listing{
		Title:               request.Title,
		OldPrice:            request.OldPrice,
		NewPrice:            request.NewPrice,
		Discount:            request.Discount,
		DiscountDescription: request.DiscountDescription,
		DietaryRestrictions: request.DietaryRestriction,
		Description:         request.Description,
		StartDate:           request.StartDate,
		StartTime:           request.StartTime,
		EndTime:             request.EndTime,
		BusinessID:          request.BusinessID,
		MultipleDays:        request.MultipleDays,
		EndDate:             request.EndDate,
		Recurring:           request.Recurring,
		RecurringDays:       request.RecurringDays,
		RecurringEndDate:    request.RecurringEndDate,
		Type:                request.ListingType,
		ImageLink:           request.ImageLink,
	}

	listingID, err := rtr.engines.AddListing(&l)
	request.ListingID = listingID

	result := listingADDResult{listingADDRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r addListingEndpoint) Validate(request interface{}) error {
	input := request.(listingADDRequest)

	if input.ListingType != shared.ListingTypeMeal && input.ListingType != shared.ListingTypeHappyHour {
		return helper.ValidationError{Message: fmt.Sprint("listing add failed, invalid 'listingType'")}
	}

	var businessFields = []string{input.Title, input.StartDate, input.StartTime, input.EndTime}
	for _, field := range businessFields {
		if strings.TrimSpace(field) == "" {
			return helper.ValidationError{Message: fmt.Sprint("listing add failed, missing mandatory fields")}
		}
	}

	if input.Discount == 0 && input.NewPrice == 0 && input.DiscountDescription == "" {
		return helper.ValidationError{Message: fmt.Sprint("listing add failed, add 'newPrice' or 'discount' or 'discountDescription' for the listing")}
	}

	if input.MultipleDays && input.Recurring {
		return helper.ValidationError{Message: fmt.Sprint("listing add failed, listing cannot be multiple days and recurring")}
	}

	if input.MultipleDays && input.EndDate == "" {
		return helper.ValidationError{Message: fmt.Sprint("listing add failed, please provide 'endDate' for multiple days lising")}
	}

	if input.Recurring && input.RecurringEndDate == "" {
		return helper.ValidationError{Message: fmt.Sprint("listing add failed, please provide 'recurringEndDate' for recurring listing")}
	} else if input.Recurring && len(input.RecurringDays) == 0 {
		return helper.ValidationError{Message: fmt.Sprint("listing add failed, please provide 'recurringDays' for recurring listing")}
	}

	return nil
}

func (r addListingEndpoint) GetPath() string {
	return "/listing/add"
}

func (r addListingEndpoint) HTTPRequest() interface{} {
	return listingADDRequest{}
}
