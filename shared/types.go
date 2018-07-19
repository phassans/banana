package shared

type (
	Listing struct {
		ListingID            int
		Title                string
		BusinessID           int
		BusinessName         string
		OldPrice             float64
		NewPrice             float64
		Discount             float64
		DietaryRestriction   []string
		Description          string
		StartDate            string
		StartTime            string
		EndTime              string
		MultipleDays         bool
		EndDate              string
		Recurring            bool
		RecurringDays        []string
		RecurringEndDate     string
		Type                 string
		ListingImage         string
		DistanceFromLocation float64
		ListingDate          string
	}

	SearchListingResult struct {
		ListingID            int      `json:"listingId"`
		ListingType          string   `json:"listingType"`
		Title                string   `json:"title"`
		BusinessID           int      `json:"businessId"`
		BusinessName         string   `json:"businessName"`
		Price                float64  `json:"price"`
		Discount             float64  `json:"discount"`
		DietaryRestriction   []string `json:"dietaryRestriction"`
		TimeLeft             int      `json:"timeLeft"`
		ListingImage         string   `json:"listingImage"`
		DistanceFromLocation float64  `json:"distanceFromLocation"`
	}

	ListingInfo struct {
		Business BusinessInfo        `json:"businessInfo"`
		Listing  SearchListingResult `json:"listing"`
	}

	ListingDate struct {
		ListingID   int
		ListingDate string
		StartTime   string
		EndTime     string
	}

	Notification struct {
		NotificationID     int
		PhoneId            string
		BusinessId         int
		Price              string
		Keywords           string
		DietaryRestriction []string
		Latitude           float64
		Longitude          float64
		Location           string
	}

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

	AddressGeo struct {
		AddressID  int
		BusinessID int
		Latitude   float64
		Longitude  float64
	}

	Business struct {
		BusinessID int    `json:"businessId"`
		Name       string `json:"name"`
		Phone      string `json:"phone"`
		Website    string `json:"website"`
	}

	BusinessAddress struct {
		Street     string `json:"street"`
		City       string `json:"city"`
		PostalCode string `json:"postalCode"`
		State      string `json:"state"`
		BusinessID int    `json:"businessID"`
	}

	BusinessInfo struct {
		Business        Business        `json:"business"`
		BusinessAddress BusinessAddress `json:"businessAddress"`
		Hours           []Bhour         `json:"businessHours"`
	}

	SortView struct {
		Listing  Listing
		Mile     float64
		Price    float64
		TimeLeft float64
	}

	CurrentLocation struct {
		Latitude  float64
		Longitude float64
	}
)
