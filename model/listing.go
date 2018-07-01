package model

import (
	"database/sql"
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
		ListingID   int
		BusinessID  int
		Title       string
		Description string
		OldPrice    float64
		NewPrice    float64
		ListingDate string
		StartTime   string
		EndTime     string
		Recurring   bool
	}

	AddressGeo struct {
		AddressID  int
		BusinessID int
		Latitude   float64
		Longitude  float64
	}
)

const insertListingSQL = "INSERT INTO listing(title, description, old_price, new_price, " +
	"listing_date, start_time, end_time, business_id, recurring, listing_create_date) " +
	"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning listing_id"

type ListingEngine interface {
	AddListing(
		title string,
		description string,
		oldPrice float64,
		newPrice float64,
		listingDate string,
		startTime string,
		endTime string,
		recurring bool,
		businessName string,
	) error
	AddListingImage(businessName string, imagePath string)
	AddListingDietaryRestrictions(listingTitle string, dietaryRestriction string)
	AddListingRecurringInfo(listingTitle string, day string, startTime string, endTime string)

	GetAllListingsInRange(
		latitude float64,
		longitude float64,
		zipCode int,
		priceFilter string,
		dietaryFilter string,
		keywords string,
	) ([]Listing, error)
}

func NewListingEngine(psql *sql.DB, logger xlog.Logger, businessEngine BusinessEngine) ListingEngine {
	return &listingEngine{psql, logger, businessEngine}
}

func (l *listingEngine) AddListing(
	title string,
	description string,
	oldPrice float64,
	newPrice float64,
	listingDate string,
	startTime string,
	endTime string,
	recurring bool,
	businessName string,
) error {
	businessID, err := l.businessEngine.GetBusinessIDFromName(businessName)
	if err != nil {
		return err
	}

	if businessID == 0 {
		return helper.BusinessDoesNotExist{BusinessName: businessName}
	}

	var listingID int

	err = l.sql.QueryRow(insertListingSQL, title, description, oldPrice, newPrice,
		listingDate, startTime, endTime, businessID, recurring, time.Now()).
		Scan(&listingID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	l.logger.Infof("successfully added a listing %s for business: %s", title, businessName)

	return nil
}

func (l *listingEngine) AddListingImage(businessName string, imagePath string) {
	return
}

func (l *listingEngine) AddListingDietaryRestrictions(listingTitle string, dietaryRestriction string) {
	return
}

func (l *listingEngine) AddListingRecurringInfo(listingTitle string, day string, startTime string, endTime string) {
	return
}

func (l *listingEngine) GetAllListingsInRange(
	latitude float64,
	longitude float64,
	zipCode int,
	priceFilter string,
	dietaryFilter string,
	keywords string,
) ([]Listing, error) {

	geoAddresses, err := l.GetAllListingsLatLon()
	if err != nil {
		return nil, err
	}

	var listings []Listing
	for _, geo := range geoAddresses {

		// find out lat, lon in range
		fromMobile := haversine.Coord{Lat: latitude, Lon: longitude}
		fromDB := haversine.Coord{Lat: geo.Latitude, Lon: geo.Longitude}
		mi, _ := haversine.Distance(fromMobile, fromDB)

		// TBD: sort on miles
		if mi > 5.0 {
			listingsFromBusinessID, err := l.GetAllListingsFromBusinessID(geo.BusinessID)
			if err != nil {
				return nil, err
			}
			listings = append(listings, listingsFromBusinessID...)
		}
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

func (l *listingEngine) GetAllListingsFromBusinessID(businessID int) ([]Listing, error) {
	currentDate := time.Now().Format("2006-01-02")
	currentTime := time.Now().Format("15:04:05.000000")

	rows, err := l.sql.Query("SELECT listing_id, business_id, title, description, old_price, new_price, "+
		"listing_date, start_time, end_time, recurring FROM listing where "+
		"business_id = $1 AND listing_date >= $2 AND end_time >= $3;", businessID, currentDate, currentTime)
	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var listings []Listing
	for rows.Next() {
		var listing Listing
		err := rows.Scan(
			&listing.ListingID,
			&listing.BusinessID,
			&listing.Title,
			&listing.Description,
			&listing.OldPrice,
			&listing.NewPrice,
			&listing.ListingDate,
			&listing.StartTime,
			&listing.EndTime,
			&listing.Recurring,
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
