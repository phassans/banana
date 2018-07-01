package controller

import (
	"context"

	"github.com/phassans/banana/model"
)

type (
	addBusinessHourRequest struct {
		BusinessID int `json:"business_id"`

		Monday                    bool   `json:"monday"`
		MondayOpenTimeSessionOne  string `json:"monday_open_time_session_one,omitempty"`
		MondayCloseTimeSessionOne string `json:"monday_close_time_session_one,omitempty"`
		MondayOpenTimeSessionTwo  string `json:"monday_open_time_session_two,omitempty"`
		MondayCloseTimeSessionTwo string `json:"monday_close_time_session_two,omitempty"`

		Tuesday                    bool   `json:"tuesday"`
		TuesdayOpenTimeSessionOne  string `json:"tuesday_open_time_session_one,omitempty"`
		TuesdayCloseTimeSessionOne string `json:"tuesday_close_time_session_one,omitempty"`
		TuesdayOpenTimeSessionTwo  string `json:"tuesday_open_time_session_two,omitempty"`
		TuesdayCloseTimeSessionTwo string `json:"tuesday_close_time_session_two,omitempty"`

		Wednesday                    bool   `json:"wednesday"`
		WednesdayOpenTimeSessionOne  string `json:"wednesday_open_time_session_one,omitempty"`
		WednesdayCloseTimeSessionOne string `json:"wednesday_close_time_session_one,omitempty"`
		WednesdayOpenTimeSessionTwo  string `json:"wednesday_open_time_session_two,omitempty"`
		WednesdayCloseTimeSessionTwo string `json:"wednesday_close_time_session_two,omitempty"`

		Thursday                    bool   `json:"thursday"`
		ThursdayOpenTimeSessionOne  string `json:"thursday_open_time_session_one,omitempty"`
		ThursdayCloseTimeSessionOne string `json:"thursday_close_time_session_one,omitempty"`
		ThursdayOpenTimeSessionTwo  string `json:"thursday_open_time_session_two,omitempty"`
		ThursdayCloseTimeSessionTwo string `json:"thursday_close_time_session_two,omitempty"`

		Friday                    bool   `json:"friday"`
		FridayOpenTimeSessionOne  string `json:"friday_open_time_session_one,omitempty"`
		FridayCloseTimeSessionOne string `json:"friday_close_time_session_one,omitempty"`
		FridayOpenTimeSessionTwo  string `json:"friday_open_time_session_two,omitempty"`
		FridayCloseTimeSessionTwo string `json:"friday_close_time_session_two,omitempty"`

		Saturday                    bool   `json:"saturday"`
		SaturdayOpenTimeSessionOne  string `json:"saturday_open_time_session_one,omitempty"`
		SaturdayCloseTimeSessionOne string `json:"saturday_close_time_session_one,omitempty"`
		SaturdayOpenTimeSessionTwo  string `json:"saturday_open_time_session_two,omitempty"`
		SaturdayCloseTimeSessionTwo string `json:"saturday_close_time_session_two,omitempty"`

		Sunday                    bool   `json:"sunday"`
		SundayOpenTimeSessionOne  string `json:"sunday_open_time_session_one,omitempty"`
		SundayCloseTimeSessionOne string `json:"sunday_close_time_session_one,omitempty"`
		SundayOpenTimeSessionTwo  string `json:"sunday_open_time_session_two,omitempty"`
		SundayCloseTimeSessionTwo string `json:"sunday_close_time_session_two,omitempty"`
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

	hoursMonday := model.HoursMonday{
		request.Monday,
		request.MondayOpenTimeSessionOne,
		request.MondayCloseTimeSessionOne,
		request.MondayOpenTimeSessionTwo,
		request.MondayCloseTimeSessionTwo,
	}

	hoursTuesday := model.HoursTuesday{
		request.Tuesday,
		request.TuesdayOpenTimeSessionOne,
		request.TuesdayCloseTimeSessionOne,
		request.TuesdayOpenTimeSessionTwo,
		request.TuesdayCloseTimeSessionTwo,
	}

	hoursWednesday := model.HoursWednesday{
		request.Wednesday,
		request.WednesdayOpenTimeSessionOne,
		request.WednesdayCloseTimeSessionOne,
		request.WednesdayOpenTimeSessionTwo,
		request.WednesdayCloseTimeSessionTwo,
	}

	hoursThursday := model.HoursThursday{
		request.Thursday,
		request.ThursdayOpenTimeSessionOne,
		request.ThursdayCloseTimeSessionOne,
		request.ThursdayOpenTimeSessionTwo,
		request.ThursdayCloseTimeSessionTwo,
	}

	hoursFriday := model.HoursFriday{
		request.Friday,
		request.FridayOpenTimeSessionOne,
		request.FridayCloseTimeSessionOne,
		request.FridayOpenTimeSessionTwo,
		request.FridayCloseTimeSessionTwo,
	}

	hoursSaturday := model.HoursSaturday{
		request.Saturday,
		request.SaturdayOpenTimeSessionOne,
		request.SaturdayCloseTimeSessionOne,
		request.SaturdayOpenTimeSessionTwo,
		request.SaturdayCloseTimeSessionTwo,
	}

	hoursSunday := model.HoursSunday{
		request.Sunday,
		request.SundayOpenTimeSessionOne,
		request.SundayCloseTimeSessionOne,
		request.SundayOpenTimeSessionTwo,
		request.SundayCloseTimeSessionTwo,
	}

	days := []model.Days{hoursMonday, hoursTuesday, hoursWednesday, hoursThursday, hoursFriday, hoursSaturday, hoursSunday}
	err := rtr.engines.AddBusinessHours(days, request.BusinessID)

	result := addBusinessHourResponse{addBusinessHourRequest: request, Error: NewAPIError(err)}
	return result, nil
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
