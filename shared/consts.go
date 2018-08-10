package shared

const (
	// See http://golang.org/pkg/time/#Parse

	// DateTimeFormat used in service
	DateTimeFormat = "2006-01-02T15:04:05Z"

	// DateFormat used to enter data
	DateFormat = "01/02/2006" //07/11/2018

	// DateFormatSQL used to get data from SQL
	DateFormatSQL = "2006-01-02" //07/11/2018

	// TimeLayout24Hour ...
	TimeLayout24Hour = "15:04:05"

	// TimeLayout12Hour ...
	TimeLayout12Hour = "3:04pm"

	// CountryID default USA
	CountryID = 1

	// SortByDistance ...
	SortByDistance = "distance"

	// SortByPrice ...
	SortByPrice = "price"

	// SortByTimeLeft ...
	SortByTimeLeft = "timeLeft"

	// SortByDateAdded ...
	SortByDateAdded = "dateAdded"

	// ListingTypeMeal ...
	ListingTypeMeal = "meal"

	// ListingTypeHappyHour ...
	ListingTypeHappyHour = "happyhour"

	// ListingEnded ...
	ListingEnded = "ended"

	// ListingScheduled ...
	ListingScheduled = "scheduled"

	// ListingActive ...
	ListingActive = "active"

	// ListingAll ...
	ListingAll = "all"

	// Undefined string
	Undefined = "undefined"
)

var (
	// DayMap of week days
	DayMap = map[string]int{"monday": 1, "tuesday": 2, "wednesday": 3, "thursday": 4, "friday": 5, "saturday": 6, "sunday": 7}

	// ListingTypes possible
	ListingTypes = []string{"", "meal", "happyhour"}
)
