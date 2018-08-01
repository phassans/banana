package listing

import (
	"bytes"
	"database/sql"
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
		"listing.recurring as recurring, listing.recurring_end_date as recurring_date, listing.listing_type as listing_type, " +
		"listing.business_id as business_id," +
		"listing.listing_id as listing_id, business.name as bname, listing_date.listing_date as listing_date"

	fromClause = "FROM listing " +
		"INNER JOIN listing_date ON listing.listing_id = listing_date.listing_id " +
		"INNER JOIN business ON listing.business_id = business.business_id"
)

func (l *listingEngine) SearchListings(
	listingTypes []string,
	future bool,
	latitude float64,
	longitude float64,
	location string,
	priceFilter float64,
	dietaryFilters []string,
	distanceFilter string,
	keywords string,
	sortBy string,
	phoneID string,
) ([]shared.SearchListingResult, error) {

	var listings []shared.Listing
	var err error

	//GetListings
	listings, err = l.GetListings(listingTypes, keywords, future)
	if err != nil {
		return nil, err
	}
	xlog.Infof("total number of listing found: %d", len(listings))

	// addDietaryRestrictionsToListings
	listings, err = l.AddDietaryRestrictionsAndImageToListings(listings)
	if err != nil {
		return nil, err
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
	xlog.Infof("search location: %v", currentLocation)

	// sort Listings based on sortBy
	sortListingEngine := NewSortListingEngine(listings, sortBy, currentLocation, l.sql)
	listings, err = sortListingEngine.SortListings()
	if err != nil {
		return nil, err
	}
	xlog.Infof("done sorting the listings. listings count: %d", len(listings))

	// filterResults
	listings, err = l.filterResults(listings, priceFilter, dietaryFilters, distanceFilter)
	xlog.Infof("applied filters. number of listings: %d", len(listings))

	if phoneID != "" {
		listings = l.tagListingsAsFavorites(listings, phoneID)
		xlog.Infof("tagging listings as favourites")
	}

	// populate searchResult
	return l.MassageAndPopulateSearchListings(listings)
}

func (l *listingEngine) filterResults(listings []shared.Listing, priceFilter float64,
	dietaryFilters []string, distanceFilter string) ([]shared.Listing, error) {

	var err error
	if priceFilter > 0.0 {
		listings, err = l.FilterByPrice(listings, priceFilter)
		xlog.Infof("applied priceFilter. count listings: %d", len(listings))
	}

	if len(dietaryFilters) > 0 {
		listings, err = l.FilterByDietaryRestrictions(listings, dietaryFilters)
		xlog.Infof("applied dietaryFilters. count listings: %d", len(listings))
	}

	if distanceFilter != "" {
		listings, err = l.FilterByDistance(listings, distanceFilter)
		xlog.Infof("applied distanceFilter. count listings: %d", len(listings))
	}

	return listings, err
}

func getListingTypeWhereClause(listingTypes []string) string {
	var listingTypesClause bytes.Buffer
	if len(listingTypes) > 0 {
		listingTypesClause.WriteString("(")
		i := 1
		for _, listingType := range listingTypes {
			if len(listingTypes) == i {
				listingTypesClause.WriteString(fmt.Sprintf("'%s')", listingType))
			} else {
				listingTypesClause.WriteString(fmt.Sprintf("'%s',", listingType))
			}
			i++
		}
	}
	return listingTypesClause.String()
}

func getWhereClause(listingTypes []string, future bool) (string, error) {
	var whereClause bytes.Buffer
	if future {
		var dateClause bytes.Buffer

		currentDate := time.Now().Format(shared.DateFormat)
		listingDate, err := time.Parse(shared.DateFormat, currentDate)
		if err != nil {
			return "", helper.DatabaseError{DBError: err.Error()}
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
	} else {
		currentDate := time.Now().Format(shared.DateFormatSQL) //"2006-01-02"
		currentTime := time.Now().Format(shared.TimeLayout24Hour)

		whereClause.WriteString(fmt.Sprintf("WHERE listing_date.listing_date = '%s' AND listing_date.end_time >= '%s'", currentDate, currentTime))
	}

	if len(listingTypes) > 0 {
		whereClause.WriteString(fmt.Sprintf(" AND listing_type in %s", getListingTypeWhereClause(listingTypes)))
	}

	return whereClause.String(), nil
}

func (l *listingEngine) GetListings(listingType []string, keywords string, future bool) ([]shared.Listing, error) {
	// determine where clause
	whereClause, err := getWhereClause(listingType, future)
	if err != nil {
		return nil, err
	}

	splitKeywordsBySpace := strings.Split(keywords, " ")
	searchKeywords := strings.Join(splitKeywordsBySpace, ",")

	var searchQuery string
	if keywords != "" {
		searchQuery = fmt.Sprintf("SELECT title, old_price, new_price, discount, description, start_date, end_date, "+
			"start_time, end_time, recurring, recurring_date, listing_type, business_id, listing_id, bname, listing_date "+
			"FROM (%s, "+
			"to_tsvector(business.name) || "+
			"to_tsvector(listing.title) || "+
			"to_tsvector(listing.description) as document "+
			"%s %s ) p_search "+
			"WHERE p_search.document @@ to_tsquery('%s');", searchSelect, fromClause, whereClause, searchKeywords)
	} else {
		searchQuery = fmt.Sprintf("%s %s %s;", searchSelect, fromClause, whereClause)
	}

	//xlog.Infof("search Query: %s", searchQuery)

	rows, err := l.sql.Query(searchQuery)
	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var listings []shared.Listing
	var sqlEndDate sql.NullString
	var sqlRecurringEndDate sql.NullString
	for rows.Next() {
		var listing shared.Listing
		err := rows.Scan(
			&listing.Title,
			&listing.OldPrice,
			&listing.NewPrice,
			&listing.Discount,
			&listing.Description,
			&listing.StartDate,
			&sqlEndDate,
			&listing.StartTime,
			&listing.EndTime,
			&listing.Recurring,
			&sqlRecurringEndDate,
			&listing.Type,
			&listing.BusinessID,
			&listing.ListingID,
			&listing.BusinessName,
			&listing.ListingDate,
		)
		if err != nil {
			return nil, helper.DatabaseError{DBError: err.Error()}
		}
		listing.EndDate = sqlEndDate.String
		listing.RecurringEndDate = sqlRecurringEndDate.String
		listings = append(listings, listing)
	}

	if err = rows.Err(); err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	return listings, nil
}

func (l *listingEngine) MassageAndPopulateSearchListings(listings []shared.Listing) ([]shared.SearchListingResult, error) {
	var listingsResult []shared.SearchListingResult
	for _, listing := range listings {
		timeLeft, err := CalculateTimeLeft(listing.ListingDate, listing.EndTime)
		if err != nil {
			return nil, err
		}

		dateTimeRange, err := DetermineDealDateTimeRange(listing.ListingDate, listing.StartTime, listing.EndTime)
		if err != nil {
			return nil, err
		}
		//xlog.Infof("dateTimeRange: %s", dateTimeRange)

		sr := shared.SearchListingResult{
			ListingID:            listing.ListingID,
			ListingType:          listing.Type,
			Title:                listing.Title,
			Description:          listing.Description,
			BusinessID:           listing.BusinessID,
			BusinessName:         listing.BusinessName,
			Price:                listing.NewPrice,
			Discount:             listing.Discount,
			DietaryRestriction:   listing.DietaryRestriction,
			TimeLeft:             timeLeft,
			ListingImage:         listing.ListingImage,
			DistanceFromLocation: listing.DistanceFromLocation,
			IsFavorite:           listing.IsFavorite,
			DateTimeRange:        dateTimeRange,
		}
		listingsResult = append(listingsResult, sr)
	}
	return listingsResult, nil
}

func CalculateTimeLeft(listingDate string, listingTime string) (int, error) {
	if listingDate == "" || listingTime == "" {
		return 0, nil
	}

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

func DetermineDealDateTimeRange(listingDate string, listingStartTime string, listingEndTime string) (string, error) {
	if listingDate == "" || listingStartTime == "" || listingEndTime == "" {
		return "", nil
	}

	// get listingDate in format
	listingDateFormatted, err := time.Parse(shared.DateFormatSQL, strings.Split(listingDate, "T")[0])
	if err != nil {
		return "", nil
	}

	// see if current day and listing day are same
	var buffer bytes.Buffer
	if time.Now().Format(shared.DateFormat) != listingDateFormatted.Format(shared.DateFormat) {
		buffer.WriteString(listingDateFormatted.Weekday().String() + ": ")
	}

	// determine startTime in format
	sTime, err := shared.GetTimeIn12HourFormat(listingStartTime)
	if err != nil {
		return "", nil
	}

	// determine endTime in format
	eTime, err := shared.GetTimeIn12HourFormat(listingEndTime)
	if err != nil {
		return "", nil
	}

	buffer.WriteString(sTime + "-" + eTime)
	return buffer.String(), nil
}
