package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
)

type (
	businessAddRequest struct {
		Name       string `json:"name"`
		Phone      string `json:"phone"`
		Website    string `json:"website,omitempty"`
		Street     string `json:"street"`
		City       string `json:"city"`
		PostalCode string `json:"postalCode"`
		State      string `json:"state"`

		Hours []struct {
			Day                 string `json:"day"`
			OpenTimeSessionOne  string `json:"open_time_session_one,omitempty"`
			CloseTimeSessionOne string `json:"close_time_session_one,omitempty"`
			OpenTimeSessionTwo  string `json:"open_time_session_two,omitempty"`
			CloseTimeSessionTwo string `json:"close_time_session_two,omitempty"`
		} `json:"hours"`

		Cuisine []string `json:"cuisine,omitempty"`

		BusinessID int `json:"businessId,omitempty"`
		AddressID  int `json:"addressId,omitempty"`
	}

	businessAddResult struct {
		businessAddRequest
		Error *APIError `json:"error,omitempty"`
	}

	createBusinessEndpoint struct{}
)

var businessAdd postEndpoint = createBusinessEndpoint{}

func (r createBusinessEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(businessAddRequest)

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	var hoursInfo []shared.Hours
	for _, day := range request.Hours {
		h := shared.Hours{day.Day, day.OpenTimeSessionOne, day.CloseTimeSessionOne, day.OpenTimeSessionTwo, day.CloseTimeSessionTwo}
		hoursInfo = append(hoursInfo, h)
	}

	businessID, addressID, err := rtr.engines.AddBusiness(
		request.Name,
		request.Phone,
		request.Website,
		request.Street,
		request.City,
		request.PostalCode,
		request.State,
		hoursInfo,
		request.Cuisine,
	)
	request.BusinessID = businessID
	request.AddressID = addressID
	result := businessAddResult{businessAddRequest: request, Error: NewAPIError(err)}

	return result, err
}

func (r createBusinessEndpoint) Validate(request interface{}) error {
	input := request.(businessAddRequest)

	var businessFields = []string{input.Name, input.Phone, input.Street, input.City, input.PostalCode, input.State}

	for _, field := range businessFields {
		if strings.TrimSpace(field) == "" {
			return helper.ValidationError{Message: fmt.Sprint("business add failed, missing mandatory fields")}
		}
	}

	if len(input.Hours) == 0 {
		return helper.ValidationError{Message: fmt.Sprint("business add failed, please add business hours")}
	}

	if len(input.Cuisine) == 0 {
		return helper.ValidationError{Message: fmt.Sprint("business add failed, please select business cuisne type")}
	}

	for _, hour := range input.Hours {
		if hour.Day == "" {
			return helper.ValidationError{Message: fmt.Sprint("business add failed, please select business days")}
		}

		// open time
		if hour.OpenTimeSessionOne == "" && hour.OpenTimeSessionTwo == "" {
			return helper.ValidationError{Message: fmt.Sprintf("business add failed, please select open time for %s", hour.Day)}
		}

		// close time
		if hour.CloseTimeSessionOne == "" && hour.CloseTimeSessionTwo == "" {
			return helper.ValidationError{Message: fmt.Sprintf("business add failed, please select close time %s", hour.Day)}
		}

		// close time
		if hour.OpenTimeSessionOne != "" && hour.CloseTimeSessionOne == "" {
			return helper.ValidationError{Message: fmt.Sprintf("business add failed, please select close time %s", hour.Day)}
		}

		// close time
		if hour.OpenTimeSessionTwo != "" && hour.CloseTimeSessionTwo == "" {
			return helper.ValidationError{Message: fmt.Sprintf("business add failed, please select close time %s", hour.Day)}
		}

		// open time
		if hour.CloseTimeSessionOne != "" && hour.OpenTimeSessionOne == "" {
			return helper.ValidationError{Message: fmt.Sprintf("business add failed, please select open time %s", hour.Day)}
		}

		// open time
		if hour.CloseTimeSessionTwo != "" && hour.OpenTimeSessionTwo == "" {
			return helper.ValidationError{Message: fmt.Sprintf("business add failed, please select open time %s", hour.Day)}
		}

	}

	return nil
}

func (r createBusinessEndpoint) GetPath() string {
	return "/business/add"
}

func (r createBusinessEndpoint) HTTPRequest() interface{} {
	return businessAddRequest{}
}
