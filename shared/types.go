package shared

type (
	BusinessUser struct {
		UserID   int    `json:"userId"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
		Phone    string `json:"phone,omitempty"`
	}

	Listing struct {
		ListingID            int      `json:"listingId,omitempty"`
		Title                string   `json:"title,omitempty"`
		BusinessID           int      `json:"businessId,omitempty"`
		BusinessName         string   `json:"businessName,omitempty"`
		OldPrice             float64  `json:"oldPrice,omitempty"`
		NewPrice             float64  `json:"newPrice,omitempty"`
		Discount             float64  `json:"discount,omitempty"`
		DietaryRestriction   []string `json:"dietaryRestriction,omitempty"`
		Description          string   `json:"description,omitempty"`
		StartDate            string   `json:"startDate,omitempty"`
		StartTime            string   `json:"startTime,omitempty"`
		EndTime              string   `json:"endTime,omitempty"`
		MultipleDays         bool     `json:"multipleDays,omitempty"`
		EndDate              string   `json:"endDate,omitempty"`
		Recurring            bool     `json:"recurring,omitempty"`
		RecurringDays        []string `json:"recurringDays,omitempty"`
		RecurringEndDate     string   `json:"recurringEndDate,omitempty"`
		Type                 string   `json:"listingType,omitempty"`
		ListingImage         string   `json:"listingImage,omitempty"`
		DistanceFromLocation float64  `json:"distanceFromLocation,omitempty"`
		ListingDate          string   `json:"listingDate,omitempty"`
		ListingStatus        string   `json:"listingStatus,omitempty"`
		TimeLeft             int      `json:"timeLeft,omitempty"`
		ImageLink            string   `json:"imageLink,omitempty"`
	}

	SearchListingResult struct {
		ListingID            int      `json:"listingId"`
		ListingType          string   `json:"listingType"`
		Title                string   `json:"title"`
		Description          string   `json:"description"`
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
		AddressID  int    `json:"addressID"`
	}

	BusinessCuisine struct {
		Cuisine []string `json:"cuisine"`
	}

	BusinessInfo struct {
		Business        Business        `json:"business"`
		BusinessAddress BusinessAddress `json:"businessAddress"`
		BusinessCuisine BusinessCuisine `json:"businessCuisine"`
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
