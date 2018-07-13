package model

import (
	"database/sql"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"bytes"

	"github.com/phassans/banana/clients"
	"github.com/phassans/banana/helper"
	"github.com/rs/xlog"
	"github.com/umahmood/haversine"
)

const (
	// See http://golang.org/pkg/time/#Parse
	dateTimeFormat = "2006-01-02T15:04:05Z"
	dateFormat     = "01/02/2006" //07/11/2018
	dateFormat1    = "2006-01-02" //07/11/2018
)

type (
	listingEngine struct {
		sql            *sql.DB
		logger         xlog.Logger
		businessEngine BusinessEngine
	}

	Listing struct {
		Title              string
		OldPrice           float64
		NewPrice           float64
		Discount           float64
		DietaryRestriction []string
		Description        string
		StartDate          string
		StartTime          string
		EndTime            string
		BusinessID         int
		MultipleDays       bool
		EndDate            string
		Recurring          bool
		RecurringDays      []string
		RecurringEndDate   string
		ListingID          int
		Type               string
		ListingImage       string
	}

	ListingInfo struct {
		BusinessInfo
		Listing
	}

	ListingDate struct {
		ListingID   int
		ListingDate string
		StartTime   string
		EndTime     string
	}
)

type ListingEngine interface {
	AddListing(listing Listing) error
	AddListingImage(businessName string, imagePath string)

	SearchListings(
		listingType string,
		future bool,
		latitude float64,
		longitude float64,
		Location string,
		priceFilter string,
		dietaryFilter string,
		keywords string,
		sortBy string,
	) ([]Listing, error)

	GetAllListings(businessID int, businessType string) ([]Listing, error)
	GetListingByID(listingID int) (Listing, error)
	GetListingInfo(listingID int) (ListingInfo, error)
}

func NewListingEngine(psql *sql.DB, logger xlog.Logger, businessEngine BusinessEngine) ListingEngine {
	return &listingEngine{psql, logger, businessEngine}
}

func (l *listingEngine) AddListing(listing Listing) error {
	business, err := l.businessEngine.GetBusinessFromID(listing.BusinessID)
	if err != nil {
		return err
	}

	if business.Name == "" {
		return helper.BusinessError{Message: fmt.Sprintf("business with id %d does not exist", listing.BusinessID)}
	}

	var listingID int
	const insertListingSQL = "INSERT INTO listing(business_id, title, old_price, new_price, discount, description," +
		"start_date, start_time, end_time, multiple_days, end_date, recurring, recurring_end_date, listing_type, listing_create_date) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) returning listing_id"

	err = l.sql.QueryRow(insertListingSQL, listing.BusinessID, listing.Title, listing.OldPrice, listing.NewPrice, listing.Discount,
		listing.Description, listing.StartDate, listing.StartTime, listing.EndTime, listing.MultipleDays, listing.EndDate,
		listing.Recurring, listing.RecurringEndDate, listing.Type, time.Now()).
		Scan(&listingID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}
	listing.ListingID = listingID

	if listing.Recurring {
		for _, day := range listing.RecurringDays {
			if err := l.AddRecurring(listingID, day); err != nil {
				return err
			}
		}
	}

	if len(listing.DietaryRestriction) > 0 {
		for _, restriction := range listing.DietaryRestriction {
			if err := l.AddDietaryRestriction(listingID, restriction); err != nil {
				return err
			}
		}
	}

	// insert into listing_date
	if err := l.AddListingDates(listing); err != nil {
		return err
	}

	l.logger.Infof("successfully added a listing %s for business: %s", listing.Title, business.Name)

	return nil
}

func (l *listingEngine) AddListingDates(listing Listing) error {
	// current listing date
	listings := []ListingDate{
		ListingDate{ListingID: listing.ListingID, ListingDate: listing.StartDate, StartTime: listing.StartTime, EndTime: listing.EndTime},
	}

	dayMap := map[string]int{"monday": 1, "tuesday": 2, "wednesday": 3, "thursday": 4, "friday": 5, "saturday": 6, "sunday": 7}

	listingDate, err := time.Parse(dateFormat, strings.Split(listing.StartDate, "T")[0])
	if err != nil {
		return err
	}

	if listing.MultipleDays {
		listingEndDate, err := time.Parse(dateFormat, strings.Split(listing.EndDate, "T")[0])
		if err != nil {
			return err
		}
		// difference b/w days
		days := listingEndDate.Sub(listingDate).Hours() / 24
		curDate := listingDate
		for i := 1; i < int(days); i++ {
			var lDate ListingDate
			nextDate := curDate.Add(time.Hour * 24)
			year, month, day := nextDate.Date()

			next := fmt.Sprintf("%d/%d/%d", int(month), day, year)
			lDate = ListingDate{ListingID: listing.ListingID, ListingDate: next, StartTime: listing.StartTime, EndTime: listing.EndTime}
			listings = append(listings, lDate)

			curDate = nextDate
		}
	}

	if listing.Recurring {
		listingRecurringDate, err := time.Parse(dateFormat, strings.Split(listing.RecurringEndDate, "T")[0])
		if err != nil {
			return err
		}
		// difference b/w days
		days := listingRecurringDate.Sub(listingDate).Hours() / 24
		curDate := listingDate
		for i := 1; i < int(days); i++ {
			var lDate ListingDate
			nextDate := curDate.Add(time.Hour * 24)
			year, month, day := nextDate.Date()
			for _, recurringDay := range listing.RecurringDays {
				if dayMap[recurringDay] == int(nextDate.Weekday()) {
					next := fmt.Sprintf("%d/%d/%d", int(month), day, year)
					lDate = ListingDate{ListingID: listing.ListingID, ListingDate: next, StartTime: listing.StartTime, EndTime: listing.EndTime}
					listings = append(listings, lDate)
				}
			}
			curDate = nextDate
		}
	}

	for _, listing := range listings {
		err := l.InsertListingDate(listing)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *listingEngine) InsertListingDate(lDate ListingDate) error {
	addListingDietRestrictionSQL := "INSERT INTO listing_date(listing_id,listing_date,start_time,end_time) " +
		"VALUES($1,$2,$3,$4);"

	_, err := l.sql.Query(addListingDietRestrictionSQL, lDate.ListingID, lDate.ListingDate, lDate.StartTime, lDate.EndTime)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	l.logger.Infof("InsertListingDate successful for listing:%d", lDate.ListingID)
	return nil
}

func (l *listingEngine) AddRecurring(listingID int, day string) error {
	addListingRecurringSQL := "INSERT INTO recurring_listing(listing_id,day) " +
		"VALUES($1,$2);"

	_, err := l.sql.Query(addListingRecurringSQL, listingID, day)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	l.logger.Infof("add recurring successful for listing:%d", listingID)
	return nil
}

func (l *listingEngine) AddDietaryRestriction(listingID int, restriction string) error {
	addListingDietRestrictionSQL := "INSERT INTO listing_dietary_restrictions(listing_id,restriction) " +
		"VALUES($1,$2);"

	_, err := l.sql.Query(addListingDietRestrictionSQL, listingID, restriction)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	l.logger.Infof("add listing_dietary_restrictions successful for listing:%d", listingID)
	return nil
}

func (l *listingEngine) AddListingImage(businessName string, imagePath string) {
	return
}

func (l *listingEngine) GetAllListings(businessID int, businessType string) ([]Listing, error) {
	getListingsQuery := "SELECT title, old_price, new_price, discount, description," +
		"start_date, end_date, start_time, end_time, recurring, listing_type, business_id, listing_id FROM listing where " +
		"business_id = $1"

	rows, err := l.sql.Query(getListingsQuery, businessID)
	if err != nil {
		return []Listing{}, helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	var listings []Listing
	for rows.Next() {
		var listing Listing
		err := rows.Scan(
			&listing.Title,
			&listing.OldPrice,
			&listing.NewPrice,
			&listing.Discount,
			&listing.Description,
			&listing.StartDate,
			&listing.EndDate,
			&listing.StartTime,
			&listing.EndTime,
			&listing.Recurring,
			&listing.Type,
			&listing.BusinessID,
			&listing.ListingID,
		)
		if err != nil {
			return []Listing{}, helper.DatabaseError{DBError: err.Error()}
		}

		// add dietary req's
		reqs, err := l.GetDietaryRestriction(listing.ListingID)
		if err != nil {
			return []Listing{}, helper.DatabaseError{DBError: err.Error()}
		}
		listing.DietaryRestriction = reqs

		// add recurring listing
		recurring, err := l.GetRecurringListing(listing.ListingID)
		if err != nil {
			return []Listing{}, helper.DatabaseError{DBError: err.Error()}
		}
		listing.RecurringDays = recurring

		listings = append(listings, listing)
	}

	if err = rows.Err(); err != nil {
		return []Listing{}, helper.DatabaseError{DBError: err.Error()}
	}

	return listings, nil
}

func (l *listingEngine) GetDietaryRestriction(listingID int) ([]string, error) {
	rows, err := l.sql.Query("SELECT restriction FROM listing_dietary_restrictions where "+
		"listing_id = $1;", listingID)
	if err != nil {
		return []string{}, helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	var dietaryReqs []string
	for rows.Next() {
		var diet string
		err := rows.Scan(
			&diet,
		)
		if err != nil {
			return []string{}, helper.DatabaseError{DBError: err.Error()}
		}
		dietaryReqs = append(dietaryReqs, diet)
	}

	if err = rows.Err(); err != nil {
		return []string{}, helper.DatabaseError{DBError: err.Error()}
	}

	return dietaryReqs, nil
}

func (l *listingEngine) GetRecurringListing(listingID int) ([]string, error) {
	rows, err := l.sql.Query("SELECT day FROM recurring_listing where "+
		"listing_id = $1;", listingID)
	if err != nil {
		return []string{}, helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	var days []string
	for rows.Next() {
		var day string
		err := rows.Scan(
			&day,
		)
		if err != nil {
			return []string{}, helper.DatabaseError{DBError: err.Error()}
		}
		days = append(days, day)
	}

	if err = rows.Err(); err != nil {
		return []string{}, helper.DatabaseError{DBError: err.Error()}
	}

	return days, nil
}

type (
	sortDistanceView struct {
		listing Listing
		mile    float64
	}

	sortPriceView struct {
		listing Listing
		price   float64
	}

	sortTimeView struct {
		listing  Listing
		timeLeft float64
	}

	CurrentLocation struct {
		Latitude  float64
		Longitude float64
	}
)

func (l *listingEngine) SearchListings(
	listingType string,
	future bool,
	latitude float64,
	longitude float64,
	location string,
	priceFilter string,
	dietaryFilter string,
	keywords string,
	sortBy string,
) ([]Listing, error) {

	var listings []Listing
	var err error
	if !future {
		//GetTodayListings
		listings, err = l.GetTodayListings(listingType)
		if err != nil {
			return nil, err
		}
	} else {
		//GetFutureListings
		listings, err = l.GetFutureListings(listingType)
		if err != nil {
			return nil, err
		}
	}

	var currentLocation CurrentLocation
	if location != "" {
		// getLatLonFromLocation
		resp, err := clients.GetLatLong(location)
		if err != nil {
			return nil, err
		}
		currentLocation = CurrentLocation{Latitude: resp.Lat, Longitude: resp.Lon}
	} else {
		currentLocation = CurrentLocation{Latitude: latitude, Longitude: longitude}
	}

	// AddDietaryRestrictionsToListings
	listings, err = l.AddDietaryRestrictionsToListings(listings)
	if err != nil {
		return nil, err
	}

	// sort Listings based on sortBy
	return l.SortListings(listings, sortBy, currentLocation, priceFilter)
}

func (l *listingEngine) AddDietaryRestrictionsToListings(listings []Listing) ([]Listing, error) {
	// get dietary restriction
	var listingsResult []Listing
	for _, listing := range listings {
		rests, err := l.GetListingsDietaryRestriction(listing.ListingID)
		if err != nil {
			return nil, err
		}
		listing.DietaryRestriction = rests
		listing.ListingImage = l.GetListingImage()
		listingsResult = append(listingsResult, listing)
	}
	return listingsResult, nil
}

func (l *listingEngine) SortListings(listings []Listing, sortingType string,
	currentLocation CurrentLocation, priceFilter string) ([]Listing, error) {

	if sortingType == "distance" || sortingType == "" {
		return l.SortListingsByDistance(listings, currentLocation)
	} else if sortingType == "price" {
		return l.SortListingsByPrice(listings, priceFilter)
	} else if sortingType == "timeLeft" {
		return l.SortListingsByTimeLeft(listings)
	}

	return nil, nil
}

func (l *listingEngine) SortListingsByTimeLeft(listings []Listing) ([]Listing, error) {
	var ll []sortTimeView
	for _, listing := range listings {

		dateTime := GetListingDateTime(listing.StartDate, listing.StartTime)
		then, err := time.Parse(dateTimeFormat, dateTime)
		if err != nil {
			return nil, nil
		}

		duration := time.Since(then)

		s := sortTimeView{listing: listing, timeLeft: duration.Seconds()}
		ll = append(ll, s)
	}

	// sort
	priceView := l.OrderListingsByTime(ll)

	// put in listing struct
	var listingsResult []Listing
	for _, view := range priceView {
		listingsResult = append(listingsResult, view.listing)
	}

	return listingsResult, nil
}

func GetListingDateTime(endDate string, endTime string) string {
	listingEndDate := strings.Split(endDate, "T")[0]
	listingEndTime := strings.Split(endTime, "T")[1]
	return fmt.Sprintf("%sT%s", listingEndDate, listingEndTime)
}

func (l *listingEngine) SortListingsByPrice(listings []Listing, priceFilter string) ([]Listing, error) {
	var ll []sortPriceView
	for _, listing := range listings {
		s := sortPriceView{listing: listing, price: listing.NewPrice}
		ll = append(ll, s)
	}

	// sort
	priceView := l.OrderListingsByPrice(ll)

	// put in listing struct
	var listingsResult []Listing
	for _, view := range priceView {
		listingsResult = append(listingsResult, view.listing)
	}

	return listingsResult, nil
}

func (l *listingEngine) SortListingsByDistance(listings []Listing, currentLocation CurrentLocation) ([]Listing, error) {
	var ll []sortDistanceView
	for _, listing := range listings {
		// get LatLon
		geo, err := l.GetListingsLatLon(listing.BusinessID)
		if err != nil {
			return nil, err
		}

		// append latLon
		fromMobile := haversine.Coord{Lat: currentLocation.Latitude, Lon: currentLocation.Longitude}
		fromDB := haversine.Coord{Lat: geo.Latitude, Lon: geo.Longitude}
		mi, _ := haversine.Distance(fromMobile, fromDB)

		fmt.Printf("business_id: %d and distance: %f \n", listing.BusinessID, mi)

		s := sortDistanceView{listing: listing, mile: mi}
		ll = append(ll, s)
	}

	// sort
	distanceView := l.OrderListingsByDistance(ll)

	// put in listing struct
	var listingsResult []Listing
	for _, view := range distanceView {
		listingsResult = append(listingsResult, view.listing)
	}

	return listingsResult, nil
}

const imageBaseURL = "http://71.198.1.192:3001"

func (l *listingEngine) GetListingImage() string {
	imgRand := random(1, 6)
	return fmt.Sprintf("%s/static/%d.jpg", imageBaseURL, imgRand)
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func (l *listingEngine) OrderListingsByTime(listings []sortTimeView) []sortTimeView {
	sort.Slice(listings, func(i, j int) bool {
		return listings[i].timeLeft > listings[j].timeLeft
	})
	return listings
}

func (l *listingEngine) OrderListingsByPrice(listings []sortPriceView) []sortPriceView {
	sort.Slice(listings, func(i, j int) bool {
		return listings[i].price < listings[j].price
	})
	return listings
}

func (l *listingEngine) OrderListingsByDistance(listings []sortDistanceView) []sortDistanceView {
	sort.Slice(listings, func(i, j int) bool {
		return listings[i].mile < listings[j].mile
	})
	return listings
}

func (l *listingEngine) GetListingsLatLon(businessID int) (AddressGeo, error) {
	rows, err := l.sql.Query("SELECT address_id, business_id, latitude, longitude  FROM address_geo WHERE business_id = $1", businessID)
	if err != nil {
		return AddressGeo{}, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	geo := AddressGeo{}
	if rows.Next() {
		err := rows.Scan(&geo.AddressID, &geo.BusinessID, &geo.Latitude, &geo.Longitude)
		if err != nil {
			return AddressGeo{}, helper.DatabaseError{DBError: err.Error()}
		}
	}

	if err = rows.Err(); err != nil {
		return AddressGeo{}, helper.DatabaseError{DBError: err.Error()}
	}

	return geo, nil
}

func (l *listingEngine) GetListingsDietaryRestriction(listingID int) ([]string, error) {
	rows, err := l.sql.Query("SELECT restriction FROM listing_dietary_restrictions WHERE listing_id = $1", listingID)
	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var rests []string
	for rows.Next() {
		var rest string
		err := rows.Scan(&rest)
		if err != nil {
			return nil, helper.DatabaseError{DBError: err.Error()}
		}
		rests = append(rests, rest)
	}

	if err = rows.Err(); err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	return rests, nil
}

func (l *listingEngine) GetTodayListings(listingType string) ([]Listing, error) {
	currentDate := time.Now().Format("2006-01-02")
	currentTime := time.Now().Format("15:04:05.000000")

	q := fmt.Sprintf("SELECT listing.title, listing.old_price, listing.new_price, listing.discount, listing.description,"+
		"listing.start_date, listing.end_date, listing.start_time, listing.end_time, listing.recurring, listing.listing_type, "+
		"listing.business_id, listing.listing_id FROM listing "+
		"INNER JOIN listing_date ON listing.listing_id = listing_date.listing_id WHERE "+
		"listing_date.listing_date = '%s' AND listing_date.end_time >= '%s' AND listing_type = '%s';", currentDate, currentTime, listingType)
	fmt.Println("Query: ", q)

	rows, err := l.sql.Query(q)
	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var listings []Listing
	for rows.Next() {
		var listing Listing
		err := rows.Scan(
			&listing.Title,
			&listing.OldPrice,
			&listing.NewPrice,
			&listing.Discount,
			&listing.Description,
			&listing.StartDate,
			&listing.EndDate,
			&listing.StartTime,
			&listing.EndTime,
			&listing.Recurring,
			&listing.Type,
			&listing.BusinessID,
			&listing.ListingID,
		)
		if err != nil {
			return nil, helper.DatabaseError{DBError: err.Error()}
		}
		listings = append(listings, listing)
	}

	if err = rows.Err(); err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	return listings, nil
}

func (l *listingEngine) GetFutureListings(listingType string) ([]Listing, error) {
	var buffer bytes.Buffer

	currentDate := time.Now().Format(dateFormat)
	listingDate, err := time.Parse(dateFormat, currentDate)

	curr := listingDate
	buffer.WriteString("(")
	for i := 0; i < 7; i++ {
		nextDate := curr.Add(time.Hour * 24)
		if i == 6 {
			buffer.WriteString(fmt.Sprintf("'%s')", strings.Split(nextDate.String(), " ")[0]))
		} else {
			buffer.WriteString(fmt.Sprintf("'%s',", strings.Split(nextDate.String(), " ")[0]))
		}
		curr = nextDate
	}

	q := fmt.Sprintf("SELECT listing.title, listing.old_price, listing.new_price, listing.discount, listing.description,"+
		"listing.start_date, listing.end_date, listing.start_time, listing.end_time, listing.recurring, listing.listing_type, "+
		"listing.business_id, listing.listing_id FROM listing "+
		"INNER JOIN listing_date ON listing.listing_id = listing_date.listing_id WHERE "+
		"listing_date IN %s AND listing_type = '%s';", buffer.String(), listingType)

	fmt.Println("Query: ", q)

	rows, err := l.sql.Query(q)
	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var listings []Listing
	for rows.Next() {
		var listing Listing
		err := rows.Scan(
			&listing.Title,
			&listing.OldPrice,
			&listing.NewPrice,
			&listing.Discount,
			&listing.Description,
			&listing.StartDate,
			&listing.EndDate,
			&listing.StartTime,
			&listing.EndTime,
			&listing.Recurring,
			&listing.Type,
			&listing.BusinessID,
			&listing.ListingID,
		)
		if err != nil {
			return nil, helper.DatabaseError{DBError: err.Error()}
		}
		listings = append(listings, listing)
	}

	if err = rows.Err(); err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	return listings, nil
}

func (l *listingEngine) GetAllListingsWithDateTime(listingType string) ([]Listing, error) {
	currentDate := time.Now().Format("2006-01-02")
	//currentTime := time.Now().Format("15:04:05.000000")

	q := fmt.Sprintf("SELECT title, old_price, new_price, discount, description,"+
		"start_date, end_date, start_time, end_time, recurring, listing_type, business_id, listing_id FROM listing WHERE "+
		"end_date >= '%s' AND listing_type = '%s';", currentDate, listingType)

	fmt.Println("qry: ", q)

	/*rows, err := l.sql.Query("SELECT title, old_price, new_price, discount, description,"+
	"start_date, end_date, start_time, end_time, recurring, listing_type, business_id, listing_id FROM listing WHERE "+
	"end_date >= $1 AND end_time >= $2 AND listing_type = $3;", currentDate, currentTime, listingType)*/

	rows, err := l.sql.Query("SELECT title, old_price, new_price, discount, description,"+
		"start_date, end_date, start_time, end_time, recurring, listing_type, business_id, listing_id FROM listing WHERE "+
		"end_date >= $1 AND listing_type = $2;", currentDate, listingType)

	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var listings []Listing
	for rows.Next() {
		var listing Listing
		err := rows.Scan(
			&listing.Title,
			&listing.OldPrice,
			&listing.NewPrice,
			&listing.Discount,
			&listing.Description,
			&listing.StartDate,
			&listing.EndDate,
			&listing.StartTime,
			&listing.EndTime,
			&listing.Recurring,
			&listing.Type,
			&listing.BusinessID,
			&listing.ListingID,
		)
		if err != nil {
			return nil, helper.DatabaseError{DBError: err.Error()}
		}
		listings = append(listings, listing)
	}

	if err = rows.Err(); err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	return listings, nil
}

func (l *listingEngine) GetAllListingsLatLon() ([]AddressGeo, error) {
	rows, err := l.sql.Query("SELECT address_id, business_id, latitude, longitude  FROM address_geo")
	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var geoAddresses []AddressGeo
	for rows.Next() {
		geo := AddressGeo{}
		err := rows.Scan(&geo.AddressID, &geo.BusinessID, &geo.Latitude, &geo.Longitude)
		if err != nil {
			return nil, helper.DatabaseError{DBError: err.Error()}
		}
		geoAddresses = append(geoAddresses, geo)
	}

	if err = rows.Err(); err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	return geoAddresses, nil
}

func (l *listingEngine) GetListingInfo(listingID int) (ListingInfo, error) {
	var listingInfo ListingInfo

	//GetListingByID
	listing, err := l.GetListingByID(listingID)
	if err != nil {
		return ListingInfo{}, err
	}
	listingInfo.Listing = listing

	//GetBusinessInfo
	businessInfo, err := l.businessEngine.GetBusinessInfo(listing.BusinessID)
	if err != nil {
		return ListingInfo{}, err
	}
	listingInfo.BusinessInfo = businessInfo
	return listingInfo, nil
}

func (l *listingEngine) GetListingByID(listingID int) (Listing, error) {
	rows, err := l.sql.Query("SELECT title, old_price, new_price, discount, description,"+
		"start_date, end_date, start_time, end_time, recurring, listing_type, business_id, listing_id FROM listing WHERE "+
		"listing_id = $1;", listingID)

	if err != nil {
		return Listing{}, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var listing Listing
	if rows.Next() {
		err := rows.Scan(
			&listing.Title,
			&listing.OldPrice,
			&listing.NewPrice,
			&listing.Discount,
			&listing.Description,
			&listing.StartDate,
			&listing.EndDate,
			&listing.StartTime,
			&listing.EndTime,
			&listing.Recurring,
			&listing.Type,
			&listing.BusinessID,
			&listing.ListingID,
		)
		if err != nil {
			return Listing{}, helper.DatabaseError{DBError: err.Error()}
		}
	}

	if err = rows.Err(); err != nil {
		return Listing{}, helper.DatabaseError{DBError: err.Error()}
	}

	return listing, nil
}
