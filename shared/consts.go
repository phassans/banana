package shared

const (
	// See http://golang.org/pkg/time/#Parse
	DateTimeFormat   = "2006-01-02T15:04:05Z"
	DateFormat       = "01/02/2006" //07/11/2018
	DateFormatSQL    = "2006-01-02" //07/11/2018
	TimeLayout24Hour = "15:04:05"
	TimeLayout12Hour = "03:04pm"
	CountryID        = 1

	SortByDistance = "distance"
	SortByPrice    = "price"
	SortByTimeLeft = "timeLeft"

	ListingTypeMeal      = "meal"
	ListingTypeHappyHour = "happyhour"
)

var (
	DayMap       = map[string]int{"monday": 1, "tuesday": 2, "wednesday": 3, "thursday": 4, "friday": 5, "saturday": 6, "sunday": 7}
	ListingTypes = []string{"", "meal", "happyhour"}
)
