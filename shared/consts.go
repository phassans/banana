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
	CountryID = 2

	// SortByDistance ...
	SortByDistance = "distance"

	// SortByPrice ...
	SortByPrice = "price"

	// SortByTimeLeft ...
	SortByTimeLeft = "timeLeft"

	// SortByDateAdded ...
	SortByDateAdded = "dateAdded"

	// SortByMostPopular ...
	SortByMostPopular = "mostpopular"

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

	// SearchToday ...
	SearchToday = "today"

	// SearchTomorrow ...
	SearchTomorrow = "tomorrow"

	// SearchThisWeek ...
	SearchThisWeek = "this week"

	// SearchNextWeek ...
	SearchNextWeek = "next week"
)

var (
	// DayMap of week days
	DayMap = map[string]int{"sunday": 0, "monday": 1, "tuesday": 2, "wednesday": 3, "thursday": 4, "friday": 5, "saturday": 6}

	// ListingTypes possible
	ListingTypes = []string{"", "meal", "happyhour"}

	imagesMap = map[string]int{
		"Asian Appetizers": 4,
		"BBQ":              2,
		"Bar snacks":       10,
		"Beer":             22,
		"Burger":           6,
		"Chinese Food":     14,
		"Cocktail":         9,
		"Coffee":           1,
		"Curry":            4,
		"Ice cream":        1,
		"Italian":          4,
		"Mediterranean":    8,
		"Mexican":          18,
		"Milk Tea Boba":    1,
		"Oysters":          11,
		"Pizza":            8,
		"Poke":             2,
		"Ramen":            4,
		"Skewers":          3,
		"Special Drinks":   10,
		"Sushi":            16,
		"Tacos":            12,
		"Thai":             5,
		"Wine":             17,
		"Wings":            8,
	}
)
