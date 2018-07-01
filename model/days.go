package model

type (
	HoursMonday struct {
		Monday                    bool   `json:"monday"`
		MondayOpenTimeSessionOne  string `json:"monday_open_time_session_one,omitempty"`
		MondayCloseTimeSessionOne string `json:"monday_close_time_session_one,omitempty"`
		MondayOpenTimeSessionTwo  string `json:"monday_open_time_session_two,omitempty"`
		MondayCloseTimeSessionTwo string `json:"monday_close_time_session_two,omitempty"`
	}

	HoursTuesday struct {
		Tuesday                    bool   `json:"tuesday"`
		TuesdayOpenTimeSessionOne  string `json:"tuesday_open_time_session_one,omitempty"`
		TuesdayCloseTimeSessionOne string `json:"tuesday_close_time_session_one,omitempty"`
		TuesdayOpenTimeSessionTwo  string `json:"tuesday_open_time_session_two,omitempty"`
		TuesdayCloseTimeSessionTwo string `json:"tuesday_close_time_session_two,omitempty"`
	}

	HoursWednesday struct {
		Wednesday                    bool   `json:"wednesday"`
		WednesdayOpenTimeSessionOne  string `json:"wednesday_open_time_session_one,omitempty"`
		WednesdayCloseTimeSessionOne string `json:"wednesday_close_time_session_one,omitempty"`
		WednesdayOpenTimeSessionTwo  string `json:"wednesday_open_time_session_two,omitempty"`
		WednesdayCloseTimeSessionTwo string `json:"wednesday_close_time_session_two,omitempty"`
	}

	HoursThursday struct {
		Thursday                    bool   `json:"thursday"`
		ThursdayOpenTimeSessionOne  string `json:"thursday_open_time_session_one"`
		ThursdayCloseTimeSessionOne string `json:"thursday_close_time_session_one"`
		ThursdayOpenTimeSessionTwo  string `json:"thursday_open_time_session_two"`
		ThursdayCloseTimeSessionTwo string `json:"thursday_close_time_session_two"`
	}

	HoursFriday struct {
		Friday                    bool   `json:"friday"`
		FridayOpenTimeSessionOne  string `json:"friday_open_time_session_one"`
		FridayCloseTimeSessionOne string `json:"friday_close_time_session_one"`
		FridayOpenTimeSessionTwo  string `json:"friday_open_time_session_two"`
		FridayCloseTimeSessionTwo string `json:"friday_close_time_session_two"`
	}

	HoursSaturday struct {
		Saturday                    bool   `json:"saturday"`
		SaturdayOpenTimeSessionOne  string `json:"saturday_open_time_session_one"`
		SaturdayCloseTimeSessionOne string `json:"saturday_close_time_session_one"`
		SaturdayOpenTimeSessionTwo  string `json:"saturday_open_time_session_two"`
		SaturdayCloseTimeSessionTwo string `json:"saturday_close_time_session_two"`
	}

	HoursSunday struct {
		Sunday                    bool   `json:"sunday"`
		SundayOpenTimeSessionOne  string `json:"sunday_open_time_session_one"`
		SundayCloseTimeSessionOne string `json:"sunday_close_time_session_one"`
		SundayOpenTimeSessionTwo  string `json:"sunday_open_time_session_two"`
		SundayCloseTimeSessionTwo string `json:"sunday_close_time_session_two"`
	}

	Days interface {
	}
)
