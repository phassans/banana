package controller

import (
	"context"

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

	var hoursInfo []shared.Hours
	for _, day := range request.Hours {
		h := shared.Hours{day.Day, day.OpenTimeSessionOne, day.CloseTimeSessionOne, day.OpenTimeSessionTwo, day.CloseTimeSessionTwo}
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
		request.Country,
		hoursInfo,
		request.Cuisine,
		request.BusinessID,
		request.AddressID,
	)
	result := businessEditResult{businessEditRequest: request, Error: NewAPIError(err)}

	return result, err
}

func (r businessEditEndpoint) Validate(request interface{}) error {
	return nil
}

func (r businessEditEndpoint) GetPath() string {
	return "/business/edit"
}

func (r businessEditEndpoint) HTTPRequest() interface{} {
	return businessEditRequest{}
}
