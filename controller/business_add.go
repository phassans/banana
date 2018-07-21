package controller

import (
	"context"

	"github.com/phassans/banana/shared"
)

type (
	businessRequest struct {
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
	}

	businessResult struct {
		businessRequest
		Error *APIError `json:"error,omitempty"`
	}

	createBusinessEndpoint struct{}
)

var businessAdd postEndpoint = createBusinessEndpoint{}

func (r createBusinessEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(businessRequest)

	var hoursInfo []shared.Hours
	for _, day := range request.Hours {
		h := shared.Hours{day.Day, day.OpenTimeSessionOne, day.CloseTimeSessionOne, day.OpenTimeSessionTwo, day.CloseTimeSessionTwo}
		hoursInfo = append(hoursInfo, h)
	}

	_, err := rtr.engines.AddBusiness(
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
	)
	result := businessResult{businessRequest: request, Error: NewAPIError(err)}

	return result, err
}

func (r createBusinessEndpoint) Validate(request interface{}) error {
	return nil
}

func (r createBusinessEndpoint) GetPath() string {
	return "/business/add"
}

func (r createBusinessEndpoint) HTTPRequest() interface{} {
	return businessRequest{}
}
