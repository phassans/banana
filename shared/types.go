package shared

import (
	"os"

	"github.com/rs/zerolog"
)

type (

	// BusinessUser info
	BusinessUser struct {
		UserID   int    `json:"userId"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
		Phone    string `json:"phone,omitempty"`
	}

	// Listing all fields
	Listing struct {
		ListingID            int           `json:"listingId"`
		Title                string        `json:"title"`
		BusinessID           int           `json:"businessId"`
		BusinessName         string        `json:"businessName"`
		OldPrice             float64       `json:"oldPrice,omitempty"`
		NewPrice             float64       `json:"newPrice,omitempty"`
		Discount             float64       `json:"discount,omitempty"`
		DiscountDescription  string        `json:"discountDescription,omitempty"`
		DietaryRestrictions  []string      `json:"dietaryRestrictions,omitempty"`
		Description          string        `json:"description,omitempty"`
		StartDate            string        `json:"startDate"`
		StartTime            string        `json:"startTime"`
		EndTime              string        `json:"endTime"`
		MultipleDays         bool          `json:"multipleDays"`
		EndDate              string        `json:"endDate,omitempty"`
		Recurring            bool          `json:"recurring"`
		RecurringDays        []string      `json:"recurringDays,omitempty"`
		RecurringEndDate     string        `json:"recurringEndDate,omitempty"`
		Type                 string        `json:"listingType"`
		ListingImage         string        `json:"listingImage,omitempty"`
		DistanceFromLocation float64       `json:"distanceFromLocation,omitempty"`
		ListingDate          string        `json:"listingDate"`
		ListingStatus        string        `json:"listingStatus,omitempty"`
		TimeLeft             int           `json:"timeLeft"`
		ImageLink            string        `json:"imageLink,omitempty"`
		Business             *BusinessInfo `json:"businessInfo,omitempty"`
		IsFavorite           bool          `json:"isFavorite"`
		DateTimeRange        string        `json:"dateTimeRange,omitempty"`
		ListingWeekDay       string        `json:"listingWeekDay,omitempty"`
		ListingDateID        int           `json:"listingDateId,omitempty"`
	}

	// SearchListingResult result of search
	SearchListingResult struct {
		ListingID            int      `json:"listingId"`
		ListingType          string   `json:"listingType"`
		Title                string   `json:"title"`
		Description          string   `json:"description"`
		BusinessID           int      `json:"businessId"`
		BusinessName         string   `json:"businessName"`
		Price                float64  `json:"price"`
		Discount             float64  `json:"discount"`
		DiscountDescription  string   `json:"discountDescription"`
		DietaryRestrictions  []string `json:"dietaryRestrictions"`
		TimeLeft             int      `json:"timeLeft"`
		ListingImage         string   `json:"listingImage"`
		DistanceFromLocation float64  `json:"distanceFromLocation"`
		IsFavorite           bool     `json:"isFavorite"`
		DateTimeRange        string   `json:"dateTimeRange"`
		ListingDateID        int      `json:"listingDateId,omitempty"`
	}

	// ListingInfo combination of Business and Listing
	ListingInfo struct {
		Business BusinessInfo        `json:"businessInfo"`
		Listing  SearchListingResult `json:"listing"`
	}

	// ListingDate fields
	ListingDate struct {
		ListingID   int
		ListingDate string
		StartTime   string
		EndTime     string
	}

	// Notification fields
	Notification struct {
		NotificationID     int
		PhoneID            string
		BusinessID         int
		Price              string
		Keywords           string
		DietaryRestriction []string
		Latitude           float64
		Longitude          float64
		Location           string
	}

	// Hours fields
	Hours struct {
		Day                 string
		OpenTimeSessionOne  string
		CloseTimeSessionOne string
		OpenTimeSessionTwo  string
		CloseTimeSessionTwo string
	}

	// BusinessHours combination of business with hours
	BusinessHours struct {
		BusinessID int
		HoursInfo  []Hours
	}

	// Bhour fields of business hour
	Bhour struct {
		Day       string `json:"day"`
		OpenTime  string `json:"openTime"`
		CloseTime string `json:"closeTime"`
	}

	// AddressGeo holds lat & lon of address
	AddressGeo struct {
		AddressID  int
		BusinessID int
		Latitude   float64
		Longitude  float64
	}

	// Business fields
	Business struct {
		BusinessID int    `json:"businessId"`
		Name       string `json:"name"`
		Phone      string `json:"phone"`
		Website    string `json:"website"`
	}

	// BusinessAddress fields of business address
	BusinessAddress struct {
		Street     string `json:"street"`
		City       string `json:"city"`
		PostalCode string `json:"postalCode"`
		State      string `json:"state"`
		BusinessID int    `json:"businessID"`
		AddressID  int    `json:"addressID"`
	}

	// BusinessCuisine list of business cuisines
	BusinessCuisine struct {
		Cuisine []string `json:"cuisine"`
	}

	// BusinessInfo all business fields
	BusinessInfo struct {
		Business        Business        `json:"business"`
		BusinessAddress BusinessAddress `json:"businessAddress"`
		BusinessCuisine BusinessCuisine `json:"businessCuisine"`
		Hours           []Bhour         `json:"businessHours,omitempty"`
		HoursFormatted  []string        `json:"businessHoursFormatted,omitempty"`
	}

	// SortView possible types
	SortView struct {
		Listing  Listing
		Mile     float64
		Price    float64
		TimeLeft float64
	}

	// CurrentLocation ...
	CurrentLocation struct {
		Latitude  float64
		Longitude float64
	}
)

var logger zerolog.Logger

func InitLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "hungryhour").
		Logger()
}

func GetLogger() zerolog.Logger {
	return logger
}
