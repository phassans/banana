package model

type (
	Hours struct {
		Day                 string
		OpenTimeSessionOne  string
		CloseTimeSessionOne string
		OpenTimeSessionTwo  string
		CloseTimeSessionTwo string
	}

	BusinessHours struct {
		BusinessID int
		HoursInfo  []Hours
	}

	Bhour struct {
		Day       string `json:"day"`
		OpenTime  string `json:"openTime"`
		CloseTime string `json:"closeTime"`
	}
)
