package listing

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/phassans/banana/clients"
	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/model/common"
	"github.com/phassans/banana/shared"
	"github.com/umahmood/haversine"
)

func (l *listingEngine) SearchListings(request shared.SearchRequest) ([]shared.SearchListingResult, error) {
	var listings []shared.Listing
	var err error

	// log search Request
	go func(req shared.SearchRequest) {
		err := l.LogSearchRequest(req)
		l.logger.Error().Msgf("LogSearchRequest returned with error: %s", err)
	}(request)

	// determine current location
	currentLocation, err := l.determineCurrentLocation(request)
	if err != nil {
		return nil, err
	}

	// GetListings
	listings, err = l.GetListings(request.ListingTypes, request.Keywords, request.Future, request.SearchDay, currentLocation)
	if err != nil {
		return nil, err
	}
	l.logger.Info().Msgf("total number of listing found: %d", len(listings))

	// populate UpVotes
	if err := l.populateUpVotes(request.PhoneID, listings); err != nil {
		return nil, err
	}

	// isLocationInRange
	/*if !isDistanceInRange(currentLocation) {
		l.logger.Error().Msgf("location not in range: %v", currentLocation)
		return []shared.SearchListingResult{}, helper.LocationError{Message: "location not in range"}
	}*/
	l.logger.Info().Msgf("search location: %v", currentLocation)

	// sort Listings based on sortBy
	sortListingEngine := NewSortListingEngine(listings, request.SortBy, currentLocation, l.sql)
	listings, err = sortListingEngine.SortListings(request.Future, request.SearchDay, request.Search, false)
	if err != nil {
		return nil, err
	}
	l.logger.Info().Msgf("done sorting the listings. listings count: %d", len(listings))

	// filterResults
	listings, err = l.filterResults(listings, request.PriceFilter, request.DietaryFilters, request.DistanceFilter)
	if err != nil {
		return nil, err
	}
	l.logger.Info().Msgf("applied filters. number of listings: %d", len(listings))

	// populate favorites
	if err := l.populateFavorites(request.PhoneID, listings); err != nil {
		return nil, err
	}

	var searchListing = make([]shared.SearchListingResult, 0)
	// getSearchListingsFromListings
	if isWeekDay(request.SearchDay) {
		searchListing, err = l.MassageAndPopulateSearchListingsWeekly(listings, false, request.SearchDay)
		if err != nil {
			return searchListing, err
		}
	} else {
		searchListing, err = l.MassageAndPopulateSearchListings(listings, false, request.SearchDay)
		if err != nil {
			return searchListing, err
		}
	}

	if request.SortBy == shared.SortByTimeLeft {
		searchListing = groupListingsBasedOnCurrentTime(searchListing)
		return searchListing, nil
	} else if request.SortBy == "" || request.SortBy == shared.SortByDistance {
		searchListing = GroupListingsOnNow(searchListing)
		return searchListing, nil
	}

	return searchListing, nil
}

func isWeekDay(sday string) bool {
	for day := range shared.DayMap {
		if day == sday {
			return true
		}
	}
	return false
}

func (l *listingEngine) populateFavorites(phoneID string, listings []shared.Listing) error {
	if phoneID != "" {
		listingIDFromFavorites, err := l.getAllFavoritesFromPhoneID(phoneID)
		if err != nil {
			return err
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

	return nil
}

func (l *listingEngine) populateUpVotes(phoneID string, listings []shared.Listing) error {
	for i := 0; i < len(listings); i++ {
		upvotes, err := l.GetUpVotes(listings[i].ListingID)
		if err != nil {
			return err
		}
		listings[i].UpVotes = upvotes

		id, err := l.GetUpVoteByPhoneID(phoneID, listings[i].ListingID)
		if err != nil {
			return err
		}

		if id > 0 {
			listings[i].IsUserVoted = true
		}
	}
	return nil
}

func (l *listingEngine) determineCurrentLocation(request shared.SearchRequest) (shared.GeoLocation, error) {
	var currentLocation shared.GeoLocation
	var resp clients.LatLong
	var err error
	if request.Location != "" && request.Latitude == 0 && request.Longitude == 0 {
		// check DB first
		currentLocation, err = l.GetGeoFromAddress(request.Location)
		if err != nil {
			l.logger.Error().Msgf("GetGeoFromAddress returned with error: %s", err)
			return currentLocation, err
		}
		if currentLocation != (shared.GeoLocation{}) {
			l.logger.Info().Msgf("GeoLocation found in DB")
			// if invalid location
			if currentLocation.Latitude == -1 && currentLocation.Longitude == -1 {
				return currentLocation, helper.LocationError{Message: "invalid location"}
			}
		} else {
			// else fetch from Google API getLatLonFromLocation
			resp, err = clients.GetLatLong(request.Location)
			if err != nil {
				return currentLocation, err
			}
			currentLocation = shared.GeoLocation{Latitude: resp.Lat, Longitude: resp.Lon}

			// cache the result in database
			go func() {
				err = l.AddGeoLocation(request.Location, currentLocation)
				if err != nil {
					l.logger.Error().Msgf("AddGeoLocation error: %s", err)
				}

			}()

			// if invalid location
			if currentLocation.Latitude == -1 && currentLocation.Longitude == -1 {
				return currentLocation, helper.LocationError{Message: "invalid location"}
			}

			l.logger.Info().Msgf("GeoLocation found in Google")
		}

		l.logger.Info().Msgf("geolocation lat: %f and lon: %f", currentLocation.Latitude, currentLocation.Longitude)
	} else {
		currentLocation = shared.GeoLocation{Latitude: request.Latitude, Longitude: request.Longitude}
	}

	return currentLocation, nil
}

func GroupListingsOnNow(listings []shared.SearchListingResult) []shared.SearchListingResult {
	var searchListingsLeft = make([]shared.SearchListingResult, 0)
	var searchListingRange = make([]shared.SearchListingResult, 0)
	for _, listing := range listings {
		if strings.Contains(listing.DateTimeRange, "left") && listing.DistanceFromLocation <= common.MaxDistanceToGroupNow {
			searchListingsLeft = append(searchListingsLeft, listing)
		} else {
			searchListingRange = append(searchListingRange, listing)
		}
	}
	return append(searchListingsLeft, searchListingRange...)
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

func (l *listingEngine) getListingTypeWhereClause(listingTypes []string) string {
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

func (l *listingEngine) getWhereClause(listingTypes []string, future bool, searchDay string, loc shared.GeoLocation) (string, error) {
	var whereClause bytes.Buffer

	timeInZone, err := getCurrentTimeInTimeZone(loc)
	if err != nil {
		return "", err
	}

	if searchDay == "today" || searchDay == "" {
		currentDate := timeInZone.Format(shared.DateFormatSQL) //"2006-01-02"
		currentTime := timeInZone.Format(shared.TimeLayout24Hour)

		whereClause.WriteString(fmt.Sprintf("WHERE listing_date.listing_date = '%s' AND listing_date.end_time >= '%s'", currentDate, currentTime))
	} else if searchDay == shared.SearchTomorrow || searchDay == shared.SearchThisWeek || searchDay == shared.SearchNextWeek {
		var dateClause bytes.Buffer

		currentDate := timeInZone.Format(shared.DateFormat)
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
			//fmt.Println("listingDate.Weekday().String()", listingDate.Weekday().String())
			endDay = 6 - shared.DayMap[strings.ToLower(listingDate.Weekday().String())]
		case shared.SearchNextWeek:
			startDay = 6 - shared.DayMap[strings.ToLower(listingDate.Weekday().String())]
			endDay = 13 - shared.DayMap[strings.ToLower(listingDate.Weekday().String())]
		}

		if searchDay == shared.SearchTomorrow || searchDay == shared.SearchThisWeek {
			curr := listingDate
			dateClause.WriteString("(")

			if endDay == 0 {
				if dateClause.String() == "(" {
					return "", nil
				}
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
	} else if searchDay != "" {
		currentDate := timeInZone.Format(shared.DateFormat)
		listingDate, err := time.Parse(shared.DateFormat, currentDate)
		if err != nil {
			return "", helper.DatabaseError{DBError: err.Error()}
		}

		curr := listingDate
		var res string
		for i := 0; i <= 7; i++ {
			nextDate := curr.Add(time.Hour * 24)
			curr = nextDate
			if strings.ToLower(nextDate.Weekday().String()) == searchDay {
				res = nextDate.String()
				break
			}
		}

		whereClause.WriteString(fmt.Sprintf("WHERE listing_date.listing_date = '%s'", strings.Split(res, " ")[0]))
	}

	if len(listingTypes) > 0 {
		whereClause.WriteString(fmt.Sprintf(" AND listing_type in %s", l.getListingTypeWhereClause(listingTypes)))
	}

	//fmt.Println(whereClause.String())
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

func (l *listingEngine) GetListings(listingType []string, keywords string, future bool, searchDay string, loc shared.GeoLocation) ([]shared.Listing, error) {
	if searchDay == shared.SearchThisWeek {
		todaysListings, err := l.getListings(listingType, keywords, future, shared.SearchToday, loc)
		if err != nil {
			return nil, err
		}

		thisWeekListings, err := l.getListings(listingType, keywords, future, shared.SearchThisWeek, loc)
		if err != nil {
			return nil, err
		}

		todaysListings = append(todaysListings, thisWeekListings...)
		return todaysListings, nil
	}
	return l.getListings(listingType, keywords, future, searchDay, loc)

}

func (l *listingEngine) getListings(listingType []string, keywords string, future bool, searchDay string, loc shared.GeoLocation) ([]shared.Listing, error) {
	// determine where clause
	whereClause, err := l.getWhereClause(listingType, future, searchDay, loc)
	if err != nil {
		return nil, err
	}

	// if no where clause nothing to do here, return
	if whereClause == "" {
		return nil, nil
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
			"start_time, end_time, multiple_days, recurring, recurring_date, listing_type, business_id, listing_id, listing_create_date, bname, latitude, longitude, listing_date_id, listing_date, path "+
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
	var sqlCreateDate sql.NullString
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
			&sqlCreateDate,
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
		listing.ListingCreateDate = sqlCreateDate.String
		listing.CurrentLocation = loc
		listings = append(listings, listing)
	}

	if err = rows.Err(); err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	return listings, nil
}

func (l *listingEngine) MassageAndPopulateSearchListings(listings []shared.Listing, isFavorite bool, searchDay string) ([]shared.SearchListingResult, error) {
	var listingsResult = make([]shared.SearchListingResult, 0)
	for _, listing := range listings {
		timeLeft, err := l.calculateTimeLeftForSearch(listing.ListingDate, listing.StartTime, listing.EndTime, listing.CurrentLocation)
		if err != nil {
			return nil, err
		}

		_, dateTimeRange, err := l.determineDealDateTimeRange(listing.ListingDate, listing.StartTime, listing.EndTime, true, timeLeft, isFavorite, listing.CurrentLocation)
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
			Upvotes:              listing.UpVotes,
			IsUserVoted:          listing.IsUserVoted,
		}
		listingsResult = append(listingsResult, sr)
	}
	return listingsResult, nil
}

func (l *listingEngine) GetDateTimeRangeForWeeklyListing(searchDay string, listingStartTime string, listingEndTime string) (string, error) {
	var buffer bytes.Buffer
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

	buffer.WriteString(strings.Title(searchDay) + " " + sTime + "-" + eTime)
	return buffer.String(), nil
}

func (l *listingEngine) MassageAndPopulateSearchListingsWeekly(listings []shared.Listing, isFavorite bool, searchDay string) ([]shared.SearchListingResult, error) {
	var listingsResult []shared.SearchListingResult
	for _, listing := range listings {

		dateTimeRange, err := l.GetDateTimeRangeForWeeklyListing(searchDay, listing.StartTime, listing.EndTime)
		if err != nil {
			return nil, nil
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
			Upvotes:              listing.UpVotes,
			IsUserVoted:          listing.IsUserVoted,
		}
		listingsResult = append(listingsResult, sr)
	}
	return listingsResult, nil
}

func (l *listingEngine) MassageAndPopulateSearchListingsFavorites(listings []shared.Listing, isFavorite bool, searchDay string) ([]shared.SearchListingResult, error) {
	var listingsResult []shared.SearchListingResult
	for _, listing := range listings {

		// get recurring info

		recDays, err := l.GetRecurringListing(listing.ListingID)
		if err != nil {
			return nil, nil
		}

		dateTimeRange, err := determineDealDateTimeRangeDetailsView(listing.StartTime, listing.EndTime, recDays)
		if err != nil {
			return nil, nil
		}

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
			Upvotes:              listing.UpVotes,
			IsUserVoted:          listing.IsUserVoted,
		}
		listingsResult = append(listingsResult, sr)
	}
	return listingsResult, nil
}

func (l *listingEngine) calculateTimeLeftForSearch(listingDate string, listingStartTime string, listingEndTime string, loc shared.GeoLocation) (int, error) {
	if listingDate == "" || listingStartTime == "" || listingEndTime == "" {
		return 0, nil
	}

	// get current date and time
	timeInZone, err := getCurrentTimeInTimeZone(loc)
	if err != nil {
		return 0, nil
	}

	currentDateTime := timeInZone.Format(shared.DateTimeFormat)
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
	hr := listingEndTimeFormatted.Hour()
	if hr >= 1 && hr <= 4 {
		listingEndTimeFormatted = listingEndTimeFormatted.Add(24 * time.Hour)
	}

	if err != nil {
		return 0, err
	}

	if !inTimeSpan(listingStartTimeFormatted, listingEndTimeFormatted, currentDateTimeFormatted) {
		return 0, nil
	}

	timeLeftInHours := listingEndTimeFormatted.Sub(currentDateTimeFormatted).Minutes()
	return int(timeLeftInHours), nil
}

func calculateTimeLeft(listingDate string, listingEndTime string, loc shared.GeoLocation) (int, error) {
	if listingDate == "" || listingEndTime == "" {
		return 0, nil
	}

	// get current date and time
	// get current date and time
	timeInZone, err := getCurrentTimeInTimeZone(loc)
	if err != nil {
		return 0, nil
	}

	currentDateTime := timeInZone.Format(shared.DateTimeFormat)
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

func (l *listingEngine) determineDealDateTimeRange(listingDate string, listingStartTime string, listingEndTime string, isSearch bool, timeLeft int, isFavorite bool, loc shared.GeoLocation) (string, string, error) {

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

		timeInZone, err := getCurrentTimeInTimeZone(loc)
		if err != nil {
			return "", "", nil
		}

		if timeInZone.Format(shared.DateFormat) != listingDateFormatted.Format(shared.DateFormat) {
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

func determineDealDateTimeRangeDetailsView(listingStartTime string, listingEndTime string, recurringDays []string) (string, error) {
	var buffer bytes.Buffer
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

	if len(recurringDays) == 7 {
		buffer.WriteString("All Days ")
	} else if len(recurringDays) == 6 || len(recurringDays) == 5 || len(recurringDays) == 4 {
		buffer.WriteString(strings.Title(recurringDays[0][0:3]) + "-" + strings.Title(recurringDays[len(recurringDays)-1][0:3]))
		buffer.WriteString(": ")
	} else if len(recurringDays) == 1 {
		buffer.WriteString(strings.Title(recurringDays[0][0:3]))
		buffer.WriteString(": ")
	} else {
		for i, day := range recurringDays {
			buffer.WriteString(strings.Title(day[0:3]))
			if i != len(recurringDays)-1 {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString(": ")
	}

	buffer.WriteString(sTime + "-" + eTime)
	return buffer.String(), nil
}

func isDistanceInRange(geoCode shared.GeoLocation) bool {
	realLocation := haversine.Coord{Lat: geoCode.Latitude, Lon: geoCode.Longitude}
	sunnyvaleLocation := haversine.Coord{Lat: float64(37.36883), Lon: float64(-122.0363496)}
	mi, _ := haversine.Distance(realLocation, sunnyvaleLocation)
	if mi < common.MaxRangeAroundSunnyvale {
		return true
	}

	return false
}

func (l *listingEngine) LogSearchRequest(searchRequest shared.SearchRequest) error {
	searchJSON, err := json.Marshal(searchRequest)
	if err != nil {
		return err
	}

	logSearchRequest := "INSERT INTO search(search_request,search_date) " +
		"VALUES($1,$2);"

	_, err = l.sql.Exec(logSearchRequest, string(searchJSON), time.Now())
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	return nil
}

func (u *listingEngine) GetUpVotes(listingID int) (int, error) {
	rows, err := u.sql.Query("SELECT upvote_id FROM upvotes WHERE listing_id = $1", listingID)
	if err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}

	return count, nil
}

func (u *listingEngine) GetUpVoteByPhoneID(phoneID string, listingID int) (int, error) {
	var upvoteID int
	rows := u.sql.QueryRow("SELECT upvote_id FROM upvotes WHERE phone_id = $1 AND listing_id = $2", phoneID, listingID)
	err := rows.Scan(&upvoteID)

	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}

	return upvoteID, nil
}
