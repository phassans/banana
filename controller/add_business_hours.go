package controller

import (
	"context"

	"github.com/phassans/banana/model"
)

type (
	addBusinessHourRequest struct {
		BusinessID int `json:"business_id"`

		Hours []struct {
			Day                 string `json:"day"`
			OpenTimeSessionOne  string `json:"open_time_session_one,omitempty"`
			CloseTimeSessionOne string `json:"close_time_session_one,omitempty"`
			OpenTimeSessionTwo  string `json:"open_time_session_two,omitempty"`
			CloseTimeSessionTwo string `json:"close_time_session_two,omitempty"`
		} `json:"hours"`
	}

	addBusinessHourResponse struct {
		addBusinessHourRequest
		Error *APIError `json:"error,omitempty"`
	}

	addBusinessHourEndpoint struct{}
)

var addBusinessHour postEndpoint = addBusinessHourEndpoint{}

func (r addBusinessHourEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(addBusinessHourRequest)
	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	var hoursInfo []model.Hours
	for _, day := range request.Hours {
		h := model.Hours{day.Day, day.OpenTimeSessionOne, day.CloseTimeSessionOne, day.OpenTimeSessionTwo, day.CloseTimeSessionTwo}
		hoursInfo = append(hoursInfo, h)
	}

	err := rtr.engines.AddBusinessHours(hoursInfo, request.BusinessID)
	result := addBusinessHourResponse{addBusinessHourRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r addBusinessHourEndpoint) Validate(request interface{}) error {
	return nil
}

func (r addBusinessHourEndpoint) GetPath() string {
	return "/business/hours/add"
}

func (r addBusinessHourEndpoint) HTTPRequest() interface{} {
	return addBusinessHourRequest{}
}
