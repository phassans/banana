package shared

const (
	// See http://golang.org/pkg/time/#Parse
	DateTimeFormat = "2006-01-02T15:04:05Z"
	DateFormat     = "01/02/2006" //07/11/2018

	CountryID = 1

	ImageBaseURL = "http://71.198.1.192:3001"
)

var (
	DayMap = map[string]int{"monday": 1, "tuesday": 2, "wednesday": 3, "thursday": 4, "friday": 5, "saturday": 6, "sunday": 7}
)
