package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
)

type (
	businessEditRequest struct {
		Name       string `json:"name"`
		Phone      string `json:"phone"`
		Website    string `json:"website"`
		Street     string `json:"street"`
		City       string `json:"city"`
		PostalCode string `json:"postalCode"`
		State      string `json:"state"`
		Country    string `json:"country"`

		Hours []struct {
			Day                 string `json:"day"`
			OpenTimeSessionOne  string `json:"open_time_session_one,omitempty"`
			CloseTimeSessionOne string `json:"close_time_session_one,omitempty"`
			OpenTimeSessionTwo  string `json:"open_time_session_two,omitempty"`
			CloseTimeSessionTwo string `json:"close_time_session_two,omitempty"`
		} `json:"hours"`

		Cuisine []string `json:"cuisine,omitempty"`

		BusinessID int `json:"businessId"`
		AddressID  int `json:"addressId"`
	}

	businessEditResult struct {
		businessEditRequest
		Error *APIError `json:"error,omitempty"`
	}

	businessEditEndpoint struct{}
)

var businessEdit postEndpoint = businessEditEndpoint{}

func (r businessEditEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(businessEditRequest)

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	var hoursInfo []shared.Hours
	for _, day := range request.Hours {
		h := shared.Hours{
			Day:                 day.Day,
			OpenTimeSessionOne:  day.OpenTimeSessionOne,
			CloseTimeSessionOne: day.CloseTimeSessionOne,
			OpenTimeSessionTwo:  day.OpenTimeSessionTwo,
			CloseTimeSessionTwo: day.CloseTimeSessionTwo,
		}
		hoursInfo = append(hoursInfo, h)
	}

	_, err := rtr.engines.BusinessEdit(
		request.Name,
		request.Phone,
		request.Website,
		request.Street,
		request.City,
		request.PostalCode,
		request.State,
		hoursInfo,
		request.Cuisine,
		request.BusinessID,
		request.AddressID,
	)
	result := businessEditResult{businessEditRequest: request, Error: NewAPIError(err)}

	return result, err
}

func (r businessEditEndpoint) Validate(request interface{}) error {
	input := request.(businessEditRequest)

	var businessFields = []string{input.Name, input.Phone, input.Street, input.City, input.PostalCode, input.State}

	for _, field := range businessFields {
		if strings.TrimSpace(field) == "" {
			return helper.ValidationError{Message: fmt.Sprint("business edit failed, missing mandatory fields")}
		}
	}

	if len(input.Hours) == 0 {
		return helper.ValidationError{Message: fmt.Sprint("business edit failed, please add business hours")}
	}

	if len(input.Cuisine) == 0 {
		return helper.ValidationError{Message: fmt.Sprint("business edit failed, please select business cuisne type")}
	}

	for _, hour := range input.Hours {
		if hour.Day == "" {
			return helper.ValidationError{Message: fmt.Sprint("business edit failed, please select business days")}
		}

		// open time
		if hour.OpenTimeSessionOne == "" && hour.OpenTimeSessionTwo == "" {
			return helper.ValidationError{Message: fmt.Sprintf("business edit failed, please select open time for %s", hour.Day)}
		}

		// close time
		if hour.CloseTimeSessionOne == "" && hour.CloseTimeSessionTwo == "" {
			return helper.ValidationError{Message: fmt.Sprintf("business edit failed, please select close time %s", hour.Day)}
		}

		// close time
		if hour.OpenTimeSessionOne != "" && hour.CloseTimeSessionOne == "" {
			return helper.ValidationError{Message: fmt.Sprintf("business edit failed, please select close time %s", hour.Day)}
		}

		// close time
		if hour.OpenTimeSessionTwo != "" && hour.CloseTimeSessionTwo == "" {
			return helper.ValidationError{Message: fmt.Sprintf("business edit failed, please select close time %s", hour.Day)}
		}

		// open time
		if hour.CloseTimeSessionOne != "" && hour.OpenTimeSessionOne == "" {
			return helper.ValidationError{Message: fmt.Sprintf("business edit failed, please select open time %s", hour.Day)}
		}

		// open time
		if hour.CloseTimeSessionTwo != "" && hour.OpenTimeSessionTwo == "" {
			return helper.ValidationError{Message: fmt.Sprintf("business edit failed, please select open time %s", hour.Day)}
		}

	}

	if input.BusinessID == 0 {
		return helper.ValidationError{Message: fmt.Sprintf("business edit failed, please provide businessId")}
	}

	if input.AddressID == 0 {
		return helper.ValidationError{Message: fmt.Sprintf("business edit failed, please provide addressId")}
	}

	return nil
}

func (r businessEditEndpoint) GetPath() string {
	return "/business/edit"
}

func (r businessEditEndpoint) HTTPRequest() interface{} {
	return businessEditRequest{}
}
