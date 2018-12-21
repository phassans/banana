package listing

import (
	"bytes"
	"database/sql"
	"fmt"

	"strings"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/model/business"
	"github.com/phassans/banana/model/common"
	"github.com/phassans/banana/shared"
	"github.com/rs/zerolog"
)

type (
	listingEngine struct {
		sql            *sql.DB
		logger         zerolog.Logger
		businessEngine business.BusinessEngine
		geoMap         map[string]shared.GeoLocation
	}

	// ListingEngine interface which holds all listing methods
	ListingEngine interface {

		// AddListing is to add a listing
		AddListing(listing *shared.Listing) (int, error)

		// SearchListings is to search for listings
		SearchListings(request shared.SearchRequest) ([]shared.SearchListingResult, error)

		// GetListingsByBusinessID returns listing based on businessID
		GetListingsByBusinessID(businessID int, businessType string) ([]shared.Listing, error)

		// GetListingByID returns listing based on ID
		GetListingByID(listingID int, businessID int, listingDateID int) (shared.Listing, error)

		// GetListingInfo returns listing info
		GetListingInfo(listingID int, listingDateID int, phoneID string) (shared.Listing, error)

		// GetListingInfo returns listing info
		UpdateListingDate(listingID int) error

		// MassageAndPopulateSearchListings to massage and populate search result
		MassageAndPopulateSearchListings([]shared.Listing, bool) ([]shared.SearchListingResult, error)

		// DeleteListing to delete the listing
		DeleteListing(listingID int) error

		// ListingEdit is to edit the listing
		ListingEdit(listing *shared.Listing) error

		// GetListingImage returns image of the listing
		GetListingImage(listingID int) (string, error)

		GetGeoFromAddress(string) (shared.GeoLocation, error)

		AddGeoLocation(string, shared.GeoLocation) error
	}
)

// NewListingEngine returns a instance of listingEngine
func NewListingEngine(psql *sql.DB, logger zerolog.Logger, businessEngine business.BusinessEngine) ListingEngine {
	//create geolocationmap
	geoMap := make(map[string]shared.GeoLocation)
	return &listingEngine{psql, logger, businessEngine, geoMap}
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
		err = rows.Scan(
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
	rows, err := l.sql.Query("SELECT day FROM listing_recurring where "+
		"listing_id = $1;", listingID)
	if err != nil {
		return []string{}, helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	var days []string
	for rows.Next() {
		var day string
		err = rows.Scan(
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

func (l *listingEngine) AddDietaryRestrictionsToListings(listings []shared.Listing) ([]shared.Listing, error) {
	// get dietary restriction
	for i := 0; i < len(listings); i++ {
		// add dietary restriction
		rests, err := l.GetListingsDietaryRestriction(listings[i].ListingID)
		if err != nil {
			return nil, err
		}
		listings[i].DietaryRestrictions = rests
	}
	return listings, nil
}

func optimizeImage(img string) string {
	logger := shared.GetLogger()
	if img != "" {
		imgParts := strings.Split(img, "/upload")
		if len(imgParts) != 2 {
			logger.Error().Msgf("image does not have two parts", img)
			return img
		}
		return fmt.Sprintf("%s/upload/q_auto,f_auto,fl_lossy%s", imgParts[0], imgParts[1])
	}
	return img
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
		err = rows.Scan(&rest)
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

func (l *listingEngine) GetListingInfo(listingID int, listingDateID int, phoneID string) (shared.Listing, error) {
	//var listingInfo shared.Listing

	//GetListingByID
	listing, err := l.GetListingByID(listingID, 0, listingDateID)
	if err != nil {
		return shared.Listing{}, err
	}

	if listing.ListingID == 0 {
		return shared.Listing{}, helper.ListingDoesNotExist{ListingID: listingID}
	}

	// add dietary req's
	reqs, err := l.GetDietaryRestriction(listing.ListingID)
	if err != nil {
		return shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
	}
	listing.DietaryRestrictions = reqs

	// get recurring info
	if listing.Recurring {
		if err != nil {
			return shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
		}
		listing.RecurringDays, err = l.GetRecurringListing(listingID)
	}

	//GetBusinessInfo
	businessInfo, err := l.businessEngine.GetBusinessInfo(listing.BusinessID)
	if err != nil {
		return shared.Listing{}, err
	}
	listing.Business = &businessInfo

	timeLeft, err := calculateTimeLeftForSearch(listing.ListingDate, listing.StartTime, listing.EndTime)
	weekday, dateTimeRange, err := determineDealDateTimeRange(listing.ListingDate, listing.StartTime, listing.EndTime, false, timeLeft, false)
	if err != nil {
		return shared.Listing{}, err
	}
	listing.TimeLeft = 0

	listing.ListingWeekDay = weekday
	listing.DateTimeRange = dateTimeRange

	if phoneID != "" {
		listing.IsFavorite = l.isFavorite(phoneID, listingID)
	}

	return messageListingsDateAndTime(listing)
}

func (l *listingEngine) UpdateListingDate(listingID int) error {
	//GetListingByID
	listing, err := l.GetListingByIDForUpdate(listingID)
	if err != nil {
		return err
	}

	if listing.ListingID == 0 {
		return helper.ListingDoesNotExist{ListingID: listingID}
	}

	// get recurring info
	if listing.Recurring {
		if err != nil {
			return helper.DatabaseError{DBError: err.Error()}
		}
		listing.RecurringDays, err = l.GetRecurringListing(listingID)
	}

	lis, err := messageListingsDateAndTime(listing)
	if err != nil {
		return err
	}

	if err := l.deleteListingDate(listingID); err != nil {
		return nil
	}

	// insert into listing_date
	if err := l.AddListingDates(&lis); err != nil {
		return err
	}

	return nil
}

func messageListingsDateAndTime(listing shared.Listing) (shared.Listing, error) {
	logger := shared.GetLogger()
	// convert StartTime
	startTime, err := shared.ConvertDBTime(listing.StartTime)
	if err != nil {
		logger.Error().Msgf("time error from DB to real for startTime: %s", startTime)
		return shared.Listing{}, fmt.Errorf("time error from DB to real for startTime: %s", startTime)
	}
	listing.StartTime = startTime

	// convert EndTime
	endTime, err := shared.ConvertDBTime(listing.EndTime)
	if err != nil {
		logger.Error().Msgf("time error from DB to real for EndTime: %s", endTime)
		return shared.Listing{}, fmt.Errorf("time error from DB to real for EndTime: %s", endTime)
	}
	listing.EndTime = endTime

	// convert startDate
	startDate, err := shared.ConvertDBDate(listing.StartDate)
	if err != nil {
		logger.Error().Msgf("date error from DB to real for startDate: %s", listing.StartDate)
		return shared.Listing{}, fmt.Errorf("date error from DB to real for startDate: %s", listing.StartDate)
	}
	listing.StartDate = startDate

	// convert EndDate
	if listing.EndDate != "" {
		endDate, err := shared.ConvertDBDate(listing.EndDate)
		if err != nil {
			logger.Error().Msgf("date error from DB to real for EndDate: %s", listing.EndDate)
			return shared.Listing{}, fmt.Errorf("date error from DB to real for EndDate: %s", listing.EndDate)
		}
		listing.EndDate = endDate
	}

	// convert RecurringEndDate
	if listing.RecurringEndDate != "" {
		recDate, err := shared.ConvertDBDate(listing.RecurringEndDate)
		if err != nil {
			logger.Error().Msgf("date error from DB to real for RecurringEndDate: %s", listing.RecurringEndDate)
			return shared.Listing{}, fmt.Errorf("date error from DB to real for RecurringEndDate: %s", listing.RecurringEndDate)
		}
		listing.RecurringEndDate = recDate
	}
	return listing, nil
}

func (l *listingEngine) GetListingByID(listingID int, businessID int, listingDateID int) (shared.Listing, error) {
	selectFields := fmt.Sprintf("%s, %s, %s, %s", common.ListingFields, common.ListingBusinessFields, common.ListingDateFields, common.ListingImageFields)

	var whereClause bytes.Buffer
	whereClause.WriteString(fmt.Sprintf(" WHERE listing.listing_id = %d", listingID))
	if businessID != 0 {
		whereClause.WriteString(fmt.Sprintf(" AND business.business_id = %d", businessID))
	}
	if listingDateID != 0 {
		whereClause.WriteString(fmt.Sprintf(" AND listing_date.listing_date_id = %d", listingDateID))
	}
	query := fmt.Sprintf("%s %s %s;", selectFields, common.FromClauseListing, whereClause.String())

	//fmt.Println("GetListingByID ", query)

	rows := l.sql.QueryRow(query)

	var listing shared.Listing
	var sqlEndDate sql.NullString
	var sqlRecurringEndDate sql.NullString
	err := rows.Scan(
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
		&listing.ListingDateID,
		&listing.ListingDate,
		&listing.ImageLink,
	)
	if err != nil {
		return shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
	}

	listing.ImageLink = optimizeImage(listing.ImageLink)
	listing.EndDate = sqlEndDate.String
	listing.RecurringEndDate = sqlRecurringEndDate.String

	return listing, nil
}

func (l *listingEngine) GetListingByIDForUpdate(listingID int) (shared.Listing, error) {
	selectFields := fmt.Sprintf("%s", common.ListingFields)

	var whereClause bytes.Buffer
	whereClause.WriteString(fmt.Sprintf(" WHERE listing.listing_id = %d", listingID))
	query := fmt.Sprintf("%s %s %s;", selectFields, "FROM listing", whereClause.String())

	//fmt.Println("GetListingByID ", query)

	rows := l.sql.QueryRow(query)

	var listing shared.Listing
	var sqlEndDate sql.NullString
	var sqlRecurringEndDate sql.NullString
	err := rows.Scan(
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
	)
	if err != nil {
		return shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
	}

	listing.ImageLink = optimizeImage(listing.ImageLink)
	listing.EndDate = sqlEndDate.String
	listing.RecurringEndDate = sqlRecurringEndDate.String

	return listing, nil
}

func (l *listingEngine) GetListingsByBusinessID(businessID int, status string) ([]shared.Listing, error) {
	selectFields := fmt.Sprintf("%s, %s, %s", common.ListingFields, common.ListingBusinessFields, common.ListingImageFields)

	var whereClause bytes.Buffer
	whereClause.WriteString(fmt.Sprintf(" WHERE listing.business_id = %d", businessID))
	query := fmt.Sprintf("%s %s %s;", selectFields, common.FromClauseListingAdmin, whereClause.String())

	//fmt.Println("GetListingsByBusinessID ", query)

	rows, err := l.sql.Query(query)
	if err != nil {
		return []shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
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
			&listing.ImageLink,
		)
		if err != nil {
			return []shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
		}
		listing.ImageLink = optimizeImage(listing.ImageLink)
		listing.EndDate = sqlEndDate.String
		listing.RecurringEndDate = sqlRecurringEndDate.String

		// add dietary req's
		var reqs []string
		reqs, err = l.GetDietaryRestriction(listing.ListingID)
		if err != nil {
			return []shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
		}
		listing.DietaryRestrictions = reqs

		// add recurring listing
		var recurring []string
		recurring, err = l.GetRecurringListing(listing.ListingID)
		if err != nil {
			return []shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
		}
		listing.RecurringDays = recurring

		// add listing status
		listing.ListingStatus = l.getListingStatus(listing)

		listings = append(listings, listing)
	}

	if err = rows.Err(); err != nil {
		return []shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
	}

	return l.filterListingBasedOnStatus(listings, status), nil
}

func (l *listingEngine) isFavorite(phoneID string, listingID int) bool {
	rows := l.sql.QueryRow("SELECT favorite_id FROM favorites where phone_id = $1 AND listing_id = $2;", phoneID, listingID)

	var favoriteID int
	err := rows.Scan(&favoriteID)
	if err != nil {
		return false
	}

	if favoriteID != 0 {
		return true
	}

	return false
}

func (l *listingEngine) getAllFavoritesFromPhoneID(phoneID string) ([]int, error) {
	rows, err := l.sql.Query("SELECT listing_id FROM favorites where phone_id = $1;", phoneID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var listingIDs []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		listingIDs = append(listingIDs, id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return listingIDs, nil
}

func (l *listingEngine) GetListingImage(listingID int) (string, error) {
	rows, err := l.sql.Query("SELECT path FROM listing_image where listing_id = $1;", listingID)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var imageLink string
	if rows.Next() {
		err = rows.Scan(&imageLink)
		if err != nil {
			return "", err
		}
	}
	if err = rows.Err(); err != nil {
		return "", err
	}
	if imageLink == "" {
		imageLink = "https://res.cloudinary.com/itshungryhour/image/upload/v1533011858/listing/NoPicAvailable.png"
	}
	return optimizeImage(imageLink), nil
}

func (l *listingEngine) GetGeoFromAddress(address string) (shared.GeoLocation, error) {
	if _, ok := l.geoMap[address]; ok {
		l.logger.Info().Msgf("found location memory map:%s", address)
		return l.geoMap[address], nil
	}

	return l.GetGeoFromAddressFromDB(address)
}

func (l *listingEngine) GetGeoFromAddressFromDB(address string) (shared.GeoLocation, error) {
	rows := l.sql.QueryRow("SELECT latitude,longitude FROM address_to_geo where address = $1;", address)

	var geo shared.GeoLocation
	err := rows.Scan(&geo.Latitude, &geo.Longitude)

	if err == sql.ErrNoRows {
		return shared.GeoLocation{}, nil
	}
	if err != nil {
		return shared.GeoLocation{}, err
	}

	// update in memory map
	l.geoMap[address] = geo

	return geo, nil
}

func (l *listingEngine) AddGeoLocation(location string, geoLocation shared.GeoLocation) error {
	addListingRecurringSQL := "INSERT INTO address_to_geo(address,latitude,longitude) " +
		"VALUES($1,$2,$3);"

	_, err := l.sql.Exec(addListingRecurringSQL, location, geoLocation.Latitude, geoLocation.Longitude)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	// update in memory map
	l.geoMap[location] = geoLocation

	l.logger.Info().Msgf("added geoLocation successful for location:%s", location)
	return nil
}
