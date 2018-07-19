package listing

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/phassans/banana/clients"
	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
	"github.com/rs/xlog"
)

const (
	searchSelect = "SELECT listing.title as title, listing.old_price as old_price, listing.new_price as new_price," +
		"listing.discount as discount, listing.description as description, listing.start_date as start_date," +
		"listing.end_date as end_date, listing.start_time as start_time, listing.end_time as end_time," +
		"listing.recurring as recurring, listing.listing_type as listing_type, listing.business_id as business_id," +
		"listing.listing_id as listing_id, business.name as bname, listing_date.listing_date as listing_date"

	fromClause = "FROM listing " +
		"INNER JOIN listing_date ON listing.listing_id = listing_date.listing_id " +
		"INNER JOIN business ON listing.business_id = business.business_id"
)

func (l *listingEngine) SearchListings(
	listingType string,
	future bool,
	latitude float64,
	longitude float64,
	location string,
	priceFilter float64,
	dietaryFilter string,
	keywords string,
	sortBy string,
) ([]shared.SearchListingResult, error) {

	var listings []shared.Listing
	var err error

	//GetListings
	listings, err = l.GetListings(listingType, keywords, future)
	if err != nil {
		return nil, err
	}

	if priceFilter > 0.0 {
		listings, err = l.FilterByPrice(listings, priceFilter)
	} else if dietaryFilter != "" {
		listings, err = l.FilterByDietaryRestrictions(listings, dietaryFilter)
	}

	var currentLocation shared.CurrentLocation
	if location != "" {
		// getLatLonFromLocation
		resp, err := clients.GetLatLong(location)
		if err != nil {
			return nil, err
		}
		currentLocation = shared.CurrentLocation{Latitude: resp.Lat, Longitude: resp.Lon}
	} else {
		currentLocation = shared.CurrentLocation{Latitude: latitude, Longitude: longitude}
	}

	// AddDietaryRestrictionsToListings
	listings, err = l.AddDietaryRestrictionsToListings(listings)
	if err != nil {
		return nil, err
	}

	// sort Listings based on sortBy
	sortListingEngine := NewSortListingEngine(listings, sortBy, currentLocation, l.sql)
	listings, err = sortListingEngine.SortListings()
	if err != nil {
		return nil, err
	}

	return l.massageAndPopulateSearchListings(listings)
}

func (l *listingEngine) GetListings(listingType string, keywords string, future bool) ([]shared.Listing, error) {
	var searchQuery string
	var whereClause bytes.Buffer
	if future {
		var dateClause bytes.Buffer

		currentDate := time.Now().Format(shared.DateFormat)
		listingDate, err := time.Parse(shared.DateFormat, currentDate)
		if err != nil {
			return nil, helper.DatabaseError{DBError: err.Error()}
		}

		curr := listingDate
		dateClause.WriteString("(")
		for i := 0; i < 7; i++ {
			nextDate := curr.Add(time.Hour * 24)
			if i == 6 {
				dateClause.WriteString(fmt.Sprintf("'%s')", strings.Split(nextDate.String(), " ")[0]))
			} else {
				dateClause.WriteString(fmt.Sprintf("'%s',", strings.Split(nextDate.String(), " ")[0]))
			}
			curr = nextDate
		}
		whereClause.WriteString(fmt.Sprintf("WHERE listing_date IN %s", dateClause.String()))
		if listingType != "" {
			whereClause.WriteString(fmt.Sprintf(" AND listing_type = '%s'", listingType))
		}
	} else {
		currentDate := time.Now().Format("2006-01-02")
		currentTime := time.Now().Format("15:04:05.000000")

		whereClause.WriteString(fmt.Sprintf("WHERE listing_date.listing_date = '%s' AND listing_date.end_time >= '%s'", currentDate, currentTime))
		if listingType != "" {
			whereClause.WriteString(fmt.Sprintf(" AND listing_type = '%s'", listingType))
		}
	}

	if keywords != "" {
		searchQuery = fmt.Sprintf("SELECT title, old_price, new_price, discount, description, start_date, end_date, "+
			"start_time, end_time, recurring, listing_type, business_id, listing_id, bname, listing_date "+
			"FROM (%s, "+
			"to_tsvector(business.name) || "+
			"to_tsvector(listing.title) || "+
			"to_tsvector(listing.description) as document "+
			"%s %s ) p_search "+
			"WHERE p_search.document @@ to_tsquery('%s');", searchSelect, fromClause, whereClause.String(), keywords)
	} else {
		searchQuery = fmt.Sprintf("%s %s %s;", searchSelect, fromClause, whereClause.String())
	}

	xlog.Infof("search Query: %s", searchQuery)

	rows, err := l.sql.Query(searchQuery)
	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var listings []shared.Listing
	for rows.Next() {
		var listing shared.Listing
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
			&listing.BusinessName,
			&listing.ListingDate,
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

func (l *listingEngine) massageAndPopulateSearchListings(listings []shared.Listing) ([]shared.SearchListingResult, error) {
	var listingsResult []shared.SearchListingResult
	for _, listing := range listings {
		timeLeft, err := calculateTimeLeft(listing.ListingDate, listing.EndTime)
		if err != nil {
			return nil, err
		}
		sr := shared.SearchListingResult{
			ListingID:            listing.ListingID,
			ListingType:          listing.Type,
			Title:                listing.Title,
			BusinessID:           listing.BusinessID,
			BusinessName:         listing.BusinessName,
			Price:                listing.NewPrice,
			Discount:             listing.Discount,
			DietaryRestriction:   listing.DietaryRestriction,
			TimeLeft:             timeLeft,
			ListingImage:         listing.ListingImage,
			DistanceFromLocation: listing.DistanceFromLocation,
		}
		listingsResult = append(listingsResult, sr)
	}
	return listingsResult, nil
}

func calculateTimeLeft(listingDate string, listingTime string) (int, error) {
	// get current date and time
	currentDateTime := time.Now().Format(shared.DateTimeFormat)
	currentDateTimeFormatted, err := time.Parse(shared.DateTimeFormat, currentDateTime)
	if err != nil {
		return 0, err
	}

	listingEndTime := GetListingDateTime(listingDate, listingTime)
	listingEndTimeFormatted, err := time.Parse(shared.DateTimeFormat, listingEndTime)
	if err != nil {
		return 0, err
	}

	timeLeftInHours := listingEndTimeFormatted.Sub(currentDateTimeFormatted).Hours()
	return int(timeLeftInHours), nil
}
