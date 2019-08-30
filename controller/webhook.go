package controller

import (
	"context"
	"os"
	"strings"

	"github.com/phassans/banana/shared"
	"github.com/rs/zerolog"
)

type (
	webhookRequest struct {
		// Business
		Name       string `json:"name"`
		Phone      string `json:"phone,omitempty"`
		Website    string `json:"website,omitempty"`
		Street     string `json:"street"`
		City       string `json:"city"`
		PostalCode string `json:"postalCode"`
		State      string `json:"state"`
		Monday     string `json:"monday"`
		Tuesday    string `json:"tuesday"`
		Wednesday  string `json:"wednesday"`
		Thursday   string `json:"thursday"`
		Friday     string `json:"friday"`
		Saturday   string `json:"saturday"`
		Sunday     string `json:"sunday"`

		// Listing

		Listings []struct {
			Title               string   `json:"title"`
			Description         string   `json:"description"`
			DiscountDescription string   `json:"discountDescription,omitempty"`
			StartDate           string   `json:"startDate"`
			RecurringEndDate    string   `json:"recurringEndDate,omitempty"`
			RecurringDays       []string `json:"recurringDays"`
			StartTime           string   `json:"startTime"`
			EndTime             string   `json:"endTime"`
			ListingID           int      `json:"listingId,omitempty"`
		} `json:"listings"`
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

	log := zerolog.New(os.Stdout)
	log.Info().Msgf("request %v", request)

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	hoursInfo := r.getBusinessHours(request)
	log.Info().Msgf("hoursInfo %v", hoursInfo)

	businessID, _, err := rtr.engines.AddBusiness(
		request.Name,
		request.Phone,
		request.Website,
		request.Street,
		request.City,
		request.PostalCode,
		request.State,
		hoursInfo,
		nil,
		1,
	)
	if err != nil {
		result := webhookResult{webhookRequest: request, Error: NewAPIError(err)}
		return result, nil
	}

	for i, listing := range request.Listings {
		// convert to lower
		for i, day := range listing.RecurringDays {
			listing.RecurringDays[i] = strings.ToLower(day)
		}

		// submit listing
		l := shared.Listing{
			Title:               listing.Title,
			DiscountDescription: listing.DiscountDescription,
			Description:         listing.Description,
			StartDate:           listing.StartDate,
			StartTime:           listing.StartTime,
			EndTime:             listing.EndTime,
			BusinessID:          businessID,
			MultipleDays:        false,
			EndDate:             listing.RecurringEndDate,
			Recurring:           true,
			RecurringDays:       listing.RecurringDays,
			RecurringEndDate:    listing.RecurringEndDate,
			Type:                "happyhour",
		}
		log.Info().Msgf("listing %d %v", i, l)

		listingID, err := rtr.engines.AddListing(&l)
		listing.ListingID = listingID
		if err != nil {
			result := webhookResult{webhookRequest: request, Error: NewAPIError(err)}
			return result, err
		}
	}

	result := webhookResult{webhookRequest: request, Error: NewAPIError(nil)}
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

func (r webhookEndpoint) getBusinessHours(request interface{}) []shared.Hours {
	input := request.(webhookRequest)

	var hoursInfo []shared.Hours
	if input.Monday != "" {
		res := r.getHours(input.Monday)
		res.Day = "monday"
		hoursInfo = append(hoursInfo, res)
	}

	if input.Tuesday != "" {
		res := r.getHours(input.Tuesday)
		res.Day = "tuesday"
		hoursInfo = append(hoursInfo, res)
	}

	if input.Wednesday != "" {
		res := r.getHours(input.Wednesday)
		res.Day = "wednesday"
		hoursInfo = append(hoursInfo, res)
	}

	if input.Thursday != "" {
		res := r.getHours(input.Thursday)
		res.Day = "thursday"
		hoursInfo = append(hoursInfo, res)
	}

	if input.Friday != "" {
		res := r.getHours(input.Friday)
		res.Day = "friday"
		hoursInfo = append(hoursInfo, res)
	}

	if input.Saturday != "" {
		res := r.getHours(input.Saturday)
		res.Day = "saturday"
		hoursInfo = append(hoursInfo, res)
	}

	if input.Sunday != "" {
		res := r.getHours(input.Sunday)
		res.Day = "sunday"
		hoursInfo = append(hoursInfo, res)
	}

	return hoursInfo
}

func (r webhookEndpoint) getHours(businessDay string) shared.Hours {
	res := shared.Hours{}
	parts := strings.Split(businessDay, ";")
	for index, part := range parts {
		if index == 0 {
			t := strings.Split(part, "-")
			res.OpenTimeSessionOne = t[1]
		} else if index == 1 {
			t := strings.Split(part, "-")
			res.CloseTimeSessionOne = t[1]
		} else if index == 2 {
			t := strings.Split(part, "-")
			res.OpenTimeSessionTwo = t[1]
		} else if index == 3 {
			t := strings.Split(part, "-")
			res.CloseTimeSessionTwo = t[1]
		}
	}

	return res
}
