package model

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/phassans/banana/helper"
	"github.com/rs/xlog"
	"github.com/umahmood/haversine"
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
		EndDate            string
		StartTime          string
		EndTime            string
		BusinessID         int
		Recurring          bool
		RecurringDays      []string
		ListingID          int
		Type               string
	}

	ListingInfo struct {
		BusinessInfo
		Listing
	}
)

type ListingEngine interface {
	AddListing(listing Listing) error
	AddListingImage(businessName string, imagePath string)

	GetAllListingsForLocation(
		listingType string,
		latitude float64,
		longitude float64,
		zipCode int,
		priceFilter string,
		dietaryFilter string,
		keywords string,
		sortBy string,
	) ([]Listing, error)
	GetAllListings(businessID int, business_type string) ([]Listing, error)
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
		"start_date, end_date, start_time, end_time, recurring, listing_type, listing_create_date) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) returning listing_id"

	err = l.sql.QueryRow(insertListingSQL, listing.BusinessID, listing.Title, listing.OldPrice, listing.NewPrice, listing.Discount,
		listing.Description, listing.StartDate, listing.EndDate, listing.StartTime, listing.EndTime, listing.Recurring,
		listing.Type, time.Now()).
		Scan(&listingID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

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

	l.logger.Infof("successfully added a listing %s for business: %s", listing.Title, business.Name)

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

type SortView struct {
	listing Listing
	mile    float64
}

func (l *listingEngine) GetAllListingsForLocation(
	listingType string,
	latitude float64,
	longitude float64,
	zipCode int,
	priceFilter string,
	dietaryFilter string,
	keywords string,
	sortBy string,
) ([]Listing, error) {

	// get all active listings
	listings, err := l.GetAllListingsWithDateTime(listingType)
	if err != nil {
		return nil, err
	}

	var ll []SortView
	for _, listing := range listings {
		// get LatLon
		geo, err := l.GetListingsLatLon(listing.BusinessID)
		if err != nil {
			return nil, err
		}

		// append latLon
		fromMobile := haversine.Coord{Lat: latitude, Lon: longitude}
		fromDB := haversine.Coord{Lat: geo.Latitude, Lon: geo.Longitude}
		mi, _ := haversine.Distance(fromMobile, fromDB)

		// get dietary restriction
		rests, err := l.GetListingsDietaryRestriction(listing.ListingID)
		if err != nil {
			return nil, err
		}
		listing.DietaryRestriction = rests

		s := SortView{listing: listing, mile: mi}
		ll = append(ll, s)
	}

	// sort
	sortView := l.OrderListings(ll)

	// put in listing struct
	var listingsResult []Listing
	for _, view := range sortView {
		listingsResult = append(listingsResult, view.listing)
	}

	// return result
	return listingsResult, nil
}

func (l *listingEngine) OrderListings(listings []SortView) []SortView {
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

func (l *listingEngine) GetAllListingsWithDateTime(listingType string) ([]Listing, error) {
	currentDate := time.Now().Format("2006-01-02")
	currentTime := time.Now().Format("15:04:05.000000")

	rows, err := l.sql.Query("SELECT title, old_price, new_price, discount, description,"+
		"start_date, end_date, start_time, end_time, recurring, listing_type, business_id, listing_id FROM listing WHERE "+
		"end_date >= $1 AND end_time >= $2 AND listing_type = $3;", currentDate, currentTime, listingType)

	/*rows, err := l.sql.Query("SELECT title, old_price, new_price, discount, description,"+
	"start_date, end_date, start_time, end_time, recurring, listing_type, business_id, listing_id FROM listing WHERE "+
	"end_date >= $1 AND listing_type = $2;", currentDate, listingType)*/

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
