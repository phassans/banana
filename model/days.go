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
)
