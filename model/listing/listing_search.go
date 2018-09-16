package listing

import (
	"bytes"
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/phassans/banana/clients"
	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/model/common"
	"github.com/phassans/banana/shared"
)

const (
/*SearchSelect = "SELECT listing.title as title, listing.old_price as old_price, listing.new_price as new_price," +
"listing.discount as discount, listing.discount_description as discount_description, listing.description as description, listing.start_date as start_date," +
"listing.end_date as end_date, listing.start_time as start_time, listing.end_time as end_time," +
"listing.multiple_days as multiple_days," +
"listing.recurring as recurring, listing.recurring_end_date as recurring_date, listing.listing_type as listing_type, " +
"listing.business_id as business_id, listing.listing_id as listing_id, " +
"business.name as bname, " +
"listing_date.listing_date_id as listing_date_id, listing_date.listing_date as listing_date, " +
"listing_image.path as path "*/
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
	searchDay string,
	sortBy string,
	phoneID string,
	search bool,
) ([]shared.SearchListingResult, error) {

	var listings []shared.Listing
	var err error

	//GetListings
	listings, err = l.GetListings(listingTypes, keywords, future, searchDay)
	if err != nil {
		return nil, err
	}
	l.logger.Info().Msgf("total number of listing found: %d", len(listings))

	// addDietaryRestrictionsToListings
	/*listings, err = l.AddDietaryRestrictionsToListings(listings)
	if err != nil {
		return nil, err
	}*/

	var currentLocation shared.GeoLocation
	var resp clients.LatLong
	if location != "" && latitude == 0 && longitude == 0 {
		// check DB first
		currentLocation, err = l.GetGeoFromAddress(location)
		if err != nil {
			l.logger.Error().Msgf("GetGeoFromAddress returned with error: %s", err)
			return nil, err
		}
		l.logger.Info().Msgf("GeoLocation found in DB")

		// else fetch from Google API
		if currentLocation == (shared.GeoLocation{}) {
			// getLatLonFromLocation
			resp, err = clients.GetLatLong(location)
			if err != nil {
				return nil, err
			}
			currentLocation = shared.GeoLocation{Latitude: resp.Lat, Longitude: resp.Lon}
			go func() {
				err = l.AddGeoLocation(location, currentLocation)
				l.logger.Error().Msgf("AddGeoLocation error: %s", err)

			}()
			l.logger.Info().Msgf("GeoLocation found in Google")
		}

		l.logger.Info().Msgf("geolocation lat: %f and lon: %f", currentLocation.Latitude, currentLocation.Longitude)
	} else {
		currentLocation = shared.GeoLocation{Latitude: latitude, Longitude: longitude}
	}
	l.logger.Info().Msgf("search location: %v", currentLocation)

	// sort Listings based on sortBy
	sortListingEngine := NewSortListingEngine(listings, sortBy, currentLocation, l.sql)
	listings, err = sortListingEngine.SortListings(future, searchDay, search, false)
	if err != nil {
		return nil, err
	}
	l.logger.Info().Msgf("done sorting the listings. listings count: %d", len(listings))

	// filterResults
	listings, err = l.filterResults(listings, priceFilter, dietaryFilters, distanceFilter)
	if err != nil {
		return nil, err
	}
	l.logger.Info().Msgf("applied filters. number of listings: %d", len(listings))

	if phoneID != "" {

		listingIDFromFavorites, err := l.getAllFavoritesFromPhoneID(phoneID)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(listings); i++ {
			listing := &listings[i]
			for _, listingIDFromFavorite := range listingIDFromFavorites {
				if listing.ListingID == listingIDFromFavorite {
					listing.IsFavorite = true
					break
				}
			}
		}

		l.logger.Info().Msgf("tagging listings as favourites")
	}

	searchListing, err := l.MassageAndPopulateSearchListings(listings, false)
	if err != nil {
		return searchListing, err
	}

	if sortBy == shared.SortByTimeLeft {
		searchListing = groupListingsBasedOnCurrentTime(searchListing)
		return searchListing, nil
	}

	return searchListing, nil
}

func groupListingsBasedOnCurrentTime(listings []shared.SearchListingResult) []shared.SearchListingResult {
	var searchListingsLeft []shared.SearchListingResult
	var searchListingRange []shared.SearchListingResult
	for _, listing := range listings {
		if strings.Contains(listing.DateTimeRange, "left") {
			searchListingsLeft = append(searchListingsLeft, listing)
		} else {
			searchListingRange = append(searchListingRange, listing)
		}
	}
	return append(searchListingsLeft, searchListingRange...)
}

func (l *listingEngine) filterResults(listings []shared.Listing, priceFilter float64,
	dietaryFilters []string, distanceFilter string) ([]shared.Listing, error) {

	var err error
	if priceFilter > 0.0 {
		listings, err = l.FilterByPrice(listings, priceFilter)
		l.logger.Info().Msgf("applied priceFilter. count listings: %d", len(listings))
	}

	if len(dietaryFilters) > 0 {
		listings, err = l.FilterByDietaryRestrictions(listings, dietaryFilters)
		l.logger.Info().Msgf("applied dietaryFilters. count listings: %d", len(listings))
	}

	if distanceFilter != "" {
		listings, err = l.FilterByDistance(listings, distanceFilter)
		l.logger.Info().Msgf("applied distanceFilter. count listings: %d", len(listings))
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

func getWhereClause(listingTypes []string, future bool, searchDay string) (string, error) {
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
		for i := 0; i < common.MaxFutureDays; i++ {
			nextDate := curr.Add(time.Hour * 24)
			if i == common.MaxFutureDays-1 {
				dateClause.WriteString(fmt.Sprintf("'%s')", strings.Split(nextDate.String(), " ")[0]))
			} else {
				dateClause.WriteString(fmt.Sprintf("'%s',", strings.Split(nextDate.String(), " ")[0]))
			}
			curr = nextDate
		}
		whereClause.WriteString(fmt.Sprintf("WHERE listing_date IN %s", dateClause.String()))
	} else if searchDay == "today" || searchDay == "" {
		currentDate := time.Now().Format(shared.DateFormatSQL) //"2006-01-02"
		currentTime := time.Now().Format(shared.TimeLayout24Hour)

		whereClause.WriteString(fmt.Sprintf("WHERE listing_date.listing_date = '%s' AND listing_date.end_time >= '%s'", currentDate, currentTime))
	} else if searchDay != "" {
		var dateClause bytes.Buffer

		currentDate := time.Now().Format(shared.DateFormat)
		listingDate, err := time.Parse(shared.DateFormat, currentDate)
		if err != nil {
			return "", helper.DatabaseError{DBError: err.Error()}
		}

		startDay := 0
		endDay := 0
		switch searchDay {
		case shared.SearchTomorrow:
			endDay = 1
		case shared.SearchThisWeek:
			fmt.Println("listingDate.Weekday().String()", listingDate.Weekday().String())
			endDay = 6 - shared.DayMap[strings.ToLower(listingDate.Weekday().String())]
		case shared.SearchNextWeek:
			startDay = 6 - shared.DayMap[strings.ToLower(listingDate.Weekday().String())]
			endDay = 13 - shared.DayMap[strings.ToLower(listingDate.Weekday().String())]
		}

		if searchDay == shared.SearchTomorrow || searchDay == shared.SearchThisWeek {
			curr := listingDate
			if searchDay == shared.SearchThisWeek {
				dateClause.WriteString(fmt.Sprintf("('%s',", strings.Split(curr.String(), " ")[0]))
			} else if searchDay == shared.SearchTomorrow {
				dateClause.WriteString("(")
			}

			if endDay == 0 {
				dateClause.WriteString(")")
			}

			for i := 0; i < endDay; i++ {
				nextDate := curr.Add(time.Hour * 24)
				if i == endDay-1 {
					dateClause.WriteString(fmt.Sprintf("'%s')", strings.Split(nextDate.String(), " ")[0]))
				} else {
					dateClause.WriteString(fmt.Sprintf("'%s',", strings.Split(nextDate.String(), " ")[0]))
				}
				curr = nextDate

			}
			whereClause.WriteString(fmt.Sprintf("WHERE listing_date IN %s", dateClause.String()))

			// added this to remove all done deals
			if searchDay == shared.SearchThisWeek {
				currentDate := time.Now().Format(shared.DateFormatSQL) //"2006-01-02"
				currentTime := time.Now().Format(shared.TimeLayout24Hour)

				whereClause.WriteString(fmt.Sprintf(" AND listing_date.listing_date = '%s' AND listing_date.end_time >= '%s'", currentDate, currentTime))
			}

		} else if searchDay == shared.SearchNextWeek {
			curr := listingDate
			dateClause.WriteString("(")
			for i := 0; i < endDay; i++ {
				nextDate := curr.Add(time.Hour * 24)
				if i >= startDay {
					if i == endDay-1 {
						dateClause.WriteString(fmt.Sprintf("'%s')", strings.Split(nextDate.String(), " ")[0]))
					} else {
						dateClause.WriteString(fmt.Sprintf("'%s',", strings.Split(nextDate.String(), " ")[0]))
					}
				}
				curr = nextDate

			}
			whereClause.WriteString(fmt.Sprintf("WHERE listing_date IN %s", dateClause.String()))
		}
	}

	if len(listingTypes) > 0 {
		whereClause.WriteString(fmt.Sprintf(" AND listing_type in %s", getListingTypeWhereClause(listingTypes)))
	}

	fmt.Println(whereClause.String())
	return whereClause.String(), nil
}

func (l *listingEngine) getKeywordsFromCategory(keywords string) ([]string, error) {
	rows, err := l.sql.Query("SELECT keyword FROM category_to_keyword where category ilike $1;", keywords)
	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var foundKeywords []string
	for rows.Next() {
		var keyword string
		err = rows.Scan(&keyword)
		if err != nil {
			return nil, helper.DatabaseError{DBError: err.Error()}
		}
		foundKeywords = append(foundKeywords, keyword)
	}

	if err = rows.Err(); err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}
	return foundKeywords, nil
}

func (l *listingEngine) GetListings(listingType []string, keywords string, future bool, searchDay string) ([]shared.Listing, error) {
	// determine where clause
	whereClause, err := getWhereClause(listingType, future, searchDay)
	if err != nil {
		return nil, err
	}

	splitKeywordsBySpace := strings.Split(keywords, " ")
	searchKeywords := strings.Join(splitKeywordsBySpace, ",")

	selectFields := fmt.Sprintf("%s, %s, %s, %s, %s", common.ListingFields, common.ListingBusinessFields, common.ListingBusinessAddressFields, common.ListingDateFields, common.ListingImageFields)

	var searchQuery string
	if keywords != "" {

		foundKeys, err := l.getKeywordsFromCategory(keywords)
		if err != nil {
			return nil, err
		}

		if len(foundKeys) > 0 {
			var noSpaceKeys []string
			for _, k := range foundKeys {
				k = strings.Replace(k, " ", " | ", -1)
				noSpaceKeys = append(noSpaceKeys, k)
			}
			searchKeywords = strings.Join(noSpaceKeys, " | ")
		}

		searchQuery = fmt.Sprintf("SELECT title, old_price, new_price, discount, discount_description, description, start_date, end_date, "+
			"start_time, end_time, multiple_days, recurring, recurring_date, listing_type, business_id, listing_id, bname, latitude, longitude, listing_date_id, listing_date, path "+
			"FROM (%s, "+
			"to_tsvector('english', business.name) || "+
			"to_tsvector('english', listing.title) || "+
			"to_tsvector('english', business_cuisine.cuisine) || "+
			"to_tsvector('english', listing.description) as document "+
			"%s %s ) p_search "+
			"WHERE p_search.document @@ to_tsquery('english', '%s');", selectFields, common.FromClauseListingWithAddress, whereClause, searchKeywords)
	} else {
		searchQuery = fmt.Sprintf("%s %s %s;", selectFields, common.FromClauseListingWithAddress, whereClause)
	}

	//l.logger.Info().Msgf("search Query: %s", searchQuery)

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
		err = rows.Scan(
			&listing.Title,
			&listing.OldPrice,
			&listing.NewPrice,
			&listing.Discount,
			&listing.DiscountDescription,
			&listing.Description,
			&listing.StartDate,
			&sqlEndDate,
			&listing.StartTime,
			&listing.EndTime,
			&listing.MultipleDays,
			&listing.Recurring,
			&sqlRecurringEndDate,
			&listing.Type,
			&listing.BusinessID,
			&listing.ListingID,
			&listing.BusinessName,
			&listing.Latitude,
			&listing.Longitude,
			&listing.ListingDateID,
			&listing.ListingDate,
			&listing.ListingImage,
		)
		if err != nil {
			return nil, helper.DatabaseError{DBError: err.Error()}
		}
		listing.ListingImage = optimizeImage(listing.ListingImage)
		listing.EndDate = sqlEndDate.String
		listing.RecurringEndDate = sqlRecurringEndDate.String
		listings = append(listings, listing)
	}

	if err = rows.Err(); err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	return listings, nil
}

func (l *listingEngine) MassageAndPopulateSearchListings(listings []shared.Listing, isFavorite bool) ([]shared.SearchListingResult, error) {
	var listingsResult []shared.SearchListingResult
	for _, listing := range listings {
		timeLeft, err := calculateTimeLeftForSearch(listing.ListingDate, listing.StartTime, listing.EndTime)
		if err != nil {
			return nil, err
		}

		_, dateTimeRange, err := determineDealDateTimeRange(listing.ListingDate, listing.StartTime, listing.EndTime, true, timeLeft, isFavorite)
		if err != nil {
			return nil, err
		}
		//l.logger.Info().Msgf("dateTimeRange: %s", dateTimeRange)

		sr := shared.SearchListingResult{
			ListingID:            listing.ListingID,
			ListingType:          listing.Type,
			Title:                listing.Title,
			Description:          listing.Description,
			BusinessID:           listing.BusinessID,
			BusinessName:         listing.BusinessName,
			Price:                listing.NewPrice,
			Discount:             listing.Discount,
			DiscountDescription:  listing.DiscountDescription,
			DietaryRestrictions:  listing.DietaryRestrictions,
			TimeLeft:             0,
			ListingImage:         listing.ListingImage,
			DistanceFromLocation: listing.DistanceFromLocation,
			IsFavorite:           listing.IsFavorite,
			DateTimeRange:        dateTimeRange,
			ListingDateID:        listing.ListingDateID,
		}
		listingsResult = append(listingsResult, sr)
	}
	return listingsResult, nil
}

func calculateTimeLeftForSearch(listingDate string, listingStartTime string, listingEndTime string) (int, error) {
	if listingDate == "" || listingStartTime == "" || listingEndTime == "" {
		return 0, nil
	}

	// get current date and time
	currentDateTime := time.Now().Format(shared.DateTimeFormat)
	currentDateTimeFormatted, err := time.Parse(shared.DateTimeFormat, currentDateTime)
	if err != nil {
		return 0, err
	}

	lStartTime := getListingDateTime(listingDate, listingStartTime)
	listingStartTimeFormatted, err := time.Parse(shared.DateTimeFormat, lStartTime)
	if err != nil {
		return 0, err
	}

	lEndTime := getListingDateTime(listingDate, listingEndTime)
	listingEndTimeFormatted, err := time.Parse(shared.DateTimeFormat, lEndTime)
	if err != nil {
		return 0, err
	}

	if !inTimeSpan(listingStartTimeFormatted, listingEndTimeFormatted, currentDateTimeFormatted) {
		return 0, nil
	}

	timeLeftInHours := listingEndTimeFormatted.Sub(currentDateTimeFormatted).Minutes()
	return int(timeLeftInHours), nil
}

func calculateTimeLeft(listingDate string, listingEndTime string) (int, error) {
	if listingDate == "" || listingEndTime == "" {
		return 0, nil
	}

	// get current date and time
	currentDateTime := time.Now().Format(shared.DateTimeFormat)
	currentDateTimeFormatted, err := time.Parse(shared.DateTimeFormat, currentDateTime)
	if err != nil {
		return 0, err
	}

	lEndTime := getListingDateTime(listingDate, listingEndTime)
	listingEndTimeFormatted, err := time.Parse(shared.DateTimeFormat, lEndTime)
	if err != nil {
		return 0, err
	}

	if strings.Contains(lEndTime, "T") {
		timeParts := strings.Split(lEndTime, "T")
		if len(timeParts) >= 1 && strings.Contains(timeParts[1], ":") {
			i, _ := strconv.ParseInt(strings.Split(timeParts[1], ":")[0], 10, 64)
			if i < 6 {
				multiplier := 24 - i
				listingEndTimeFormatted = listingEndTimeFormatted.Add(time.Hour * time.Duration(multiplier))
			}
		}
	}

	timeLeftInHours := listingEndTimeFormatted.Sub(currentDateTimeFormatted).Minutes()
	return int(timeLeftInHours), nil
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func determineDealDateTimeRange(listingDate string, listingStartTime string, listingEndTime string, isSearch bool, timeLeft int, isFavorite bool) (string, string, error) {

	if listingDate == "" || listingStartTime == "" || listingEndTime == "" {
		return "", "", nil
	}

	// get listingDate in format
	listingDateFormatted, err := time.Parse(shared.DateFormatSQL, strings.Split(listingDate, "T")[0])
	if err != nil {
		return "", "", nil
	}
	if timeLeft == 0 {
		// see if current day and listing day are same
		var buffer bytes.Buffer
		if !isFavorite {
			buffer.WriteString(fmt.Sprintf("%s %d, ", listingDateFormatted.Month().String()[0:3], listingDateFormatted.Day()))
		}

		if time.Now().Format(shared.DateFormat) != listingDateFormatted.Format(shared.DateFormat) {
			if isSearch {
				buffer.WriteString(listingDateFormatted.Weekday().String()[0:3] + ": ")
			} else {
				buffer.WriteString(listingDateFormatted.Weekday().String() + ": ")
			}
		} else {
			buffer.WriteString("Today: ")
		}

		// determine startTime in format
		sTime, err := shared.GetTimeIn12HourFormat(listingStartTime)
		if err != nil {
			return "", "", nil
		}

		// determine endTime in format
		eTime, err := shared.GetTimeIn12HourFormat(listingEndTime)
		if err != nil {
			return "", "", nil
		}

		buffer.WriteString(sTime + "-" + eTime)
		return listingDateFormatted.Weekday().String(), buffer.String(), nil
	} else {
		weekDayToday := listingDateFormatted.Weekday().String()
		if timeLeft < 50 {
			return weekDayToday, fmt.Sprintf("%d mins left", timeLeft), nil
		}

		resMod := math.Mod(float64(timeLeft), 60)
		if resMod < 30 {
			res := float64(timeLeft) / float64(60)
			hrs := int(math.Floor(res))
			if hrs == 1 {
				return weekDayToday, fmt.Sprintf("%d hour left", int(math.Floor(res))), nil
			}
			return weekDayToday, fmt.Sprintf("%d hours left", int(math.Floor(res))), nil
		}

		res := float64(timeLeft) / float64(60)
		hrs := int(math.Ceil(res))
		if hrs == 1 {
			return weekDayToday, fmt.Sprintf("%d hour left", int(math.Ceil(res))), nil
		}
		return weekDayToday, fmt.Sprintf("%d hours left", int(math.Ceil(res))), nil
	}
}
