package listing

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/phassans/banana/clients"
	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
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
	if !future {
		//GetTodayListings
		listings, err = l.GetTodayListings(listingType, keywords)
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

	if priceFilter > 0.0 {
		listings, err = l.FilterByPrice(listings, priceFilter)
	} else if dietaryFilter != "" {
		listings, err = l.FilterByDietaryRestrictions(listings, dietaryFilter)
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
	listings, err = l.SortListings(listings, sortBy, currentLocation)
	if err != nil {
		return nil, err
	}

	return l.massageAndPopulateSearchListings(listings)
}

func (l *listingEngine) massageAndPopulateSearchListings(listings []shared.Listing) ([]shared.SearchListingResult, error) {
	// get current date and time
	currentDateTime := time.Now().Format(shared.DateTimeFormat)
	currentDateTimeFormatted, err := time.Parse(shared.DateTimeFormat, currentDateTime)
	if err != nil {
		return nil, err
	}

	var listingsResult []shared.SearchListingResult
	for _, listing := range listings {
		listingEndTime := GetListingDateTime(listing.ListingDate, listing.EndTime)
		listingEndTimeFormatted, err := time.Parse(shared.DateTimeFormat, listingEndTime)
		if err != nil {
			return nil, err
		}

		timeLeftInHours := listingEndTimeFormatted.Sub(currentDateTimeFormatted).Hours()
		sr := shared.SearchListingResult{
			ListingID:            listing.ListingID,
			ListingType:          listing.Type,
			Title:                listing.Title,
			BusinessID:           listing.BusinessID,
			BusinessName:         listing.BusinessName,
			Price:                listing.NewPrice,
			Discount:             listing.Discount,
			DietaryRestriction:   listing.DietaryRestriction,
			TimeLeft:             int(timeLeftInHours),
			ListingImage:         listing.ListingImage,
			DistanceFromLocation: listing.DistanceFromLocation,
		}
		listingsResult = append(listingsResult, sr)
	}

	return listingsResult, nil
}

func (l *listingEngine) GetTodayListings(listingType string, keywords string) ([]shared.Listing, error) {
	currentDate := time.Now().Format("2006-01-02")
	currentTime := time.Now().Format("15:04:05.000000")

	var query string

	if keywords != "" {
		query = fmt.Sprintf(`SELECT title, old_price, new_price, discount, description, start_date, end_date, start_time, end_time, recurring, listing_type, business_id, listing_id, bname, listing_date
FROM (SELECT listing.title as title, listing.old_price as old_price, listing.new_price as new_price,
    listing.discount as discount, listing.description as description, listing.start_date as start_date,
    listing.end_date as end_date, listing.start_time as start_time, listing.end_time as end_time,
    listing.recurring as recurring, listing.listing_type as listing_type, listing.business_id as business_id,
    listing.listing_id as listing_id, business.name as bname, listing_date.listing_date as listing_date,
	to_tsvector(business.name) ||
    to_tsvector(listing.title) ||
    to_tsvector(listing.description) as document
    FROM listing
    INNER JOIN listing_date ON listing.listing_id = listing_date.listing_id
    INNER JOIN business ON listing.business_id = business.business_id
    WHERE listing_date.listing_date = '%s' AND listing_date.end_time >= '%s' AND listing_type = '%s'
) p_search
WHERE p_search.document @@ to_tsquery('%s');`, currentDate, currentTime, listingType, keywords)
	} else if listingType == "" {
		query = fmt.Sprintf("SELECT listing.title, listing.old_price, listing.new_price, listing.discount, listing.description,"+
			"listing.start_date, listing.end_date, listing.start_time, listing.end_time, listing.recurring, listing.listing_type, "+
			"listing.business_id, listing.listing_id, business.name, listing_date.listing_date FROM listing "+
			"INNER JOIN listing_date ON listing.listing_id = listing_date.listing_id "+
			"INNER JOIN business ON listing.business_id = business.business_id WHERE "+
			"listing_date.listing_date = '%s' AND listing_date.end_time >= '%s';", currentDate, currentTime)
	} else {
		query = fmt.Sprintf("SELECT listing.title, listing.old_price, listing.new_price, listing.discount, listing.description,"+
			"listing.start_date, listing.end_date, listing.start_time, listing.end_time, listing.recurring, listing.listing_type, "+
			"listing.business_id, listing.listing_id, business.name, listing_date.listing_date FROM listing "+
			"INNER JOIN listing_date ON listing.listing_id = listing_date.listing_id "+
			"INNER JOIN business ON listing.business_id = business.business_id WHERE "+
			"listing_date.listing_date = '%s' AND listing_date.end_time >= '%s' AND listing_type = '%s';", currentDate, currentTime, listingType)
	}

	fmt.Println("Query: ", query)

	rows, err := l.sql.Query(query)
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

func (l *listingEngine) GetFutureListings(listingType string) ([]shared.Listing, error) {
	var buffer bytes.Buffer

	currentDate := time.Now().Format(shared.DateFormat)
	listingDate, err := time.Parse(shared.DateFormat, currentDate)

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

	var query string
	if listingType == "" {
		query = fmt.Sprintf("SELECT listing.title, listing.old_price, listing.new_price, listing.discount, listing.description,"+
			"listing.start_date, listing.end_date, listing.start_time, listing.end_time, listing.recurring, listing.listing_type, "+
			"listing.business_id, listing.listing_id, business.name, listing_date.listing_date FROM listing "+
			"INNER JOIN listing_date ON listing.listing_id = listing_date.listing_id "+
			"INNER JOIN business ON listing.business_id = business.business_id WHERE "+
			"listing_date IN %s;", buffer.String())
	} else {
		query = fmt.Sprintf("SELECT listing.title, listing.old_price, listing.new_price, listing.discount, listing.description,"+
			"listing.start_date, listing.end_date, listing.start_time, listing.end_time, listing.recurring, listing.listing_type, "+
			"listing.business_id, listing.listing_id, business.name, listing_date.listing_date FROM listing "+
			"INNER JOIN listing_date ON listing.listing_id = listing_date.listing_id "+
			"INNER JOIN business ON listing.business_id = business.business_id WHERE "+
			"listing_date IN %s AND listing_type = '%s';", buffer.String(), listingType)
	}

	fmt.Println("Query: ", query)

	rows, err := l.sql.Query(query)
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
