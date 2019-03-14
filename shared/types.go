package shared

import (
	"os"
	"time"

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
		ListingCreateDate    string        `json:"listingCreateDate"`
		ListingStatus        string        `json:"listingStatus,omitempty"`
		TimeLeft             int           `json:"timeLeft"`
		ImageLink            string        `json:"imageLink,omitempty"`
		Business             *BusinessInfo `json:"businessInfo,omitempty"`
		IsFavorite           bool          `json:"isFavorite"`
		DateTimeRange        string        `json:"dateTimeRange,omitempty"`
		ListingWeekDay       string        `json:"listingWeekDay,omitempty"`
		ListingDateID        int           `json:"listingDateId,omitempty"`
		Favorite             *Favorite     `json:"favorite,omitempty"`
		Latitude             float64       `json:"latitude,omitempty"`
		Longitude            float64       `json:"longitude,omitempty"`
		UpVotes              int           `json:"upvotes"`
		IsUserVoted          bool          `json:"isUserUpVoted"`
		SubmittedBy          string        `json:"submittedBy"`
		CurrentLocation      GeoLocation   `json:"geoLocation,omitempty"`
	}

	// SearchListingResult result of search
	SearchListingResult struct {
		ListingID                  int      `json:"listingId"`
		ListingType                string   `json:"listingType"`
		Title                      string   `json:"title"`
		Description                string   `json:"description"`
		BusinessID                 int      `json:"businessId"`
		BusinessName               string   `json:"businessName"`
		Price                      float64  `json:"price"`
		Discount                   float64  `json:"discount"`
		DiscountDescription        string   `json:"discountDescription"`
		DietaryRestrictions        []string `json:"dietaryRestrictions"`
		TimeLeft                   int      `json:"timeLeft"`
		ListingImage               string   `json:"listingImage"`
		DistanceFromLocation       float64  `json:"distanceFromLocation"`
		DistanceFromLocationString string   `json:"distanceFromLocationString"`
		IsFavorite                 bool     `json:"isFavorite"`
		DateTimeRange              string   `json:"dateTimeRange"`
		ListingDateID              int      `json:"listingDateId,omitempty"`
		Upvotes                    int      `json:"upvotes"`
		IsUserVoted                bool     `json:"isUserUpVoted"`
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

	notificationADDRequest struct {
		PhoneID        string   `json:"phoneId"`
		Latitude       float64  `json:"latitude,omitempty"`
		Longitude      float64  `json:"longitude,omitempty"`
		Location       string   `json:"location,omitempty"`
		PriceFilter    float64  `json:"priceFilter,omitempty"`
		DietaryFilters []string `json:"dietaryFilters,omitempty"`
		DistanceFilter string   `json:"distanceFilter,omitempty"`
		Keywords       string   `json:"keywords,omitempty"`
	}

	// Notification fields
	Notification struct {
		NotificationID   int
		NotificationName string `json:"notificationName,omitempty"`
		PhoneID          string
		Latitude         float64  `json:"latitude,omitempty"`
		Longitude        float64  `json:"longitude,omitempty"`
		Location         string   `json:"location,omitempty"`
		PriceFilter      string   `json:"priceFilter,omitempty"`
		DietaryFilters   []string `json:"dietaryFilters,omitempty"`
		DistanceFilter   string   `json:"distanceFilter,omitempty"`
		Keywords         string   `json:"keywords,omitempty"`
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

	BusinessD struct {
		BusinessID int    `json:"businessId"`
		Name       string `json:"name"`
		Phone      string `json:"phone"`
		Website    string `json:"website"`
		Street     string `json:"street"`
		City       string `json:"city"`
		PostalCode string `json:"postalCode"`
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
		Street     string  `json:"street"`
		City       string  `json:"city"`
		PostalCode string  `json:"postalCode"`
		State      string  `json:"state"`
		BusinessID int     `json:"businessID"`
		AddressID  int     `json:"addressID"`
		Latitude   float64 `json:"latitude"`
		Longitude  float64 `json:"longitude"`
	}

	// BusinessCuisine list of business cuisines
	BusinessCuisine struct {
		Cuisine []string `json:"cuisine"`
	}

	// BusinessInfo all business fields
	BusinessInfo struct {
		Business        Business        `json:"business"`
		BusinessAddress BusinessAddress `json:"businessAddress"`
		BusinessCuisine BusinessCuisine `json:"businessCuisine,omitempty"`
		Hours           []Bhour         `json:"businessHours,omitempty"`
		HoursFormatted  []string        `json:"businessHoursFormatted,omitempty"`
	}

	// SortView possible types
	SortView struct {
		Listing     Listing
		Mile        float64
		Price       float64
		TimeLeft    float64
		ListingDate time.Time
		UpVotes     int
	}

	// CurrentLocation ...
	GeoLocation struct {
		Latitude  float64
		Longitude float64
	}

	// Favorite ...
	Favorite struct {
		FavoriteID      int    `json:"favoriteId,omitempty"`
		ListingID       int    `json:"listingId"`
		ListingDateID   int    `json:"listingDateId,omitempty"`
		FavoriteAddDate string `json:"favoriteAddDate,omitempty"`
	}

	Preference struct {
		PreferenceID int    `json:"preferenceId"`
		PhoneID      string `json:"phoneId"`
		Cuisine      string `json:"cuisine"`
	}

	SearchRequest struct {
		Future         bool     `json:"future"`
		Search         bool     `json:"search,omitempty"`
		ListingTypes   []string `json:"listingTypes,omitempty"`
		Latitude       float64  `json:"latitude,omitempty"`
		Longitude      float64  `json:"longitude,omitempty"`
		Location       string   `json:"location,omitempty"`
		PriceFilter    float64  `json:"priceFilter,omitempty"`
		DietaryFilters []string `json:"dietaryFilters,omitempty"`
		DistanceFilter string   `json:"distanceFilter,omitempty"`
		Keywords       string   `json:"keywords,omitempty"`
		SortBy         string   `json:"sortBy,omitempty"`
		SearchDay      string   `json:"searchDay,omitempty"`
		PhoneID        string   `json:"phoneId"`
	}
)

var logger zerolog.Logger

// InitLogger is to initialize a logger
func InitLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "hungryhour").
		Logger()
}

// GetLogger is to get Logger
func GetLogger() zerolog.Logger {
	return logger
}
