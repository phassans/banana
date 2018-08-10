package listing

import (
	"bytes"
	"database/sql"
	"fmt"

	"strings"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/model/business"
	"github.com/phassans/banana/shared"
	"github.com/rs/zerolog"
)

type (
	listingEngine struct {
		sql            *sql.DB
		logger         zerolog.Logger
		businessEngine business.BusinessEngine
	}

	// ListingEngine interface which holds all listing methods
	ListingEngine interface {

		// AddListing is to add a listing
		AddListing(listing *shared.Listing) (int, error)

		// SearchListings is to search for listings
		SearchListings(
			listingType []string,
			future bool,
			latitude float64,
			longitude float64,
			Location string,
			priceFilter float64,
			dietaryFilter []string,
			distanceFilter string,
			keywords string,
			sortBy string,
			phoneID string,
		) ([]shared.SearchListingResult, error)

		// GetListingsByBusinessID returns listing based on businessID
		GetListingsByBusinessID(businessID int, businessType string) ([]shared.Listing, error)

		// GetListingByID returns listing based on ID
		GetListingByID(listingID int, businessID int, listingDateID int) (shared.Listing, error)

		// GetListingInfo returns listing info
		GetListingInfo(listingID int, listingDateID int, phoneID string) (shared.Listing, error)

		// MassageAndPopulateSearchListings to massage and populate search result
		MassageAndPopulateSearchListings([]shared.Listing) ([]shared.SearchListingResult, error)

		// DeleteListing to delete the listing
		DeleteListing(listingID int) error

		// ListingEdit is to edit the listing
		ListingEdit(listing *shared.Listing) error
	}
)

// NewListingEngine returns a instance of listingEngine
func NewListingEngine(psql *sql.DB, logger zerolog.Logger, businessEngine business.BusinessEngine) ListingEngine {
	return &listingEngine{psql, logger, businessEngine}
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
	var listingsResult []shared.Listing
	for _, listing := range listings {
		// add dietary restriction
		rests, err := l.GetListingsDietaryRestriction(listing.ListingID)
		if err != nil {
			return nil, err
		}
		listing.DietaryRestrictions = rests
		listingsResult = append(listingsResult, listing)
	}
	return listingsResult, nil
}

func optimizeImage(img string) string {
	imgParts := strings.Split(img, "/upload")
	return fmt.Sprintf("%s/upload/w_400,c_fill%s", imgParts[0], imgParts[1])
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

	//GetBusinessInfo
	businessInfo, err := l.businessEngine.GetBusinessInfo(listing.BusinessID)
	if err != nil {
		return shared.Listing{}, err
	}
	listing.Business = &businessInfo

	weekday, dateTimeRange, err := determineDealDateTimeRange(listing.ListingDate, listing.StartTime, listing.EndTime, false)
	if err != nil {
		return shared.Listing{}, err
	}
	listing.ListingWeekDay = weekday
	listing.DateTimeRange = dateTimeRange

	timeLeft, err := calculateTimeLeftForSearch(listing.ListingDate, listing.StartTime, listing.EndTime)
	listing.TimeLeft = timeLeft

	if phoneID != "" {
		listing.IsFavorite = l.isFavorite(phoneID, listingID)
	}

	return listing, nil
}

func (l *listingEngine) GetListingByID(listingID int, businessID int, listingDateID int) (shared.Listing, error) {

	var whereClause bytes.Buffer
	whereClause.WriteString(fmt.Sprintf(" WHERE listing.listing_id = %d", listingID))
	if businessID != 0 {
		whereClause.WriteString(fmt.Sprintf(" AND business.business_id = %d", businessID))
	}
	if listingDateID != 0 {
		whereClause.WriteString(fmt.Sprintf(" AND listing_date.listing_date_id = %d", listingDateID))
	}
	query := fmt.Sprintf("%s %s %s;", SearchSelect, fromClause, whereClause.String())

	fmt.Println("GetListingByID ", query)

	rows, err := l.sql.Query(query)

	if err != nil {
		return shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var listing shared.Listing
	var sqlEndDate sql.NullString
	var sqlRecurringEndDate sql.NullString
	if rows.Next() {
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
			&listing.ListingDateID,
			&listing.ListingDate,
			&listing.ListingImage,
		)
		if err != nil {
			return shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
		}
	}
	listing.StartDate = sqlEndDate.String
	listing.RecurringEndDate = sqlRecurringEndDate.String

	if err = rows.Err(); err != nil {
		return shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
	}

	return listing, nil
}

func (l *listingEngine) GetListingsByBusinessID(businessID int, status string) ([]shared.Listing, error) {
	getListingsQuery := "SELECT listing.title as title, listing.old_price as old_price, listing.new_price as new_price," +
		"listing.discount as discount, listing.discount_description as discount_description, listing.description as description, listing.start_date as start_date," +
		"listing.end_date as end_date, listing.start_time as start_time, listing.end_time as end_time," +
		"listing.multiple_days as multiple_days," +
		"listing.recurring as recurring, listing.recurring_end_date as recurring_date, listing.listing_type as listing_type, " +
		"business.business_id as business_id, listing.listing_id as listing_id, business.name as bname " +
		"FROM listing " +
		"INNER JOIN business ON listing.business_id = business.business_id " +
		"WHERE " +
		"listing.business_id = $1"

	rows, err := l.sql.Query(getListingsQuery, businessID)
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
		)
		if err != nil {
			return []shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
		}
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

func (l *listingEngine) tagListingsAsFavorites(listings []shared.Listing, phoneID string) []shared.Listing {
	// get dietary restriction
	var result []shared.Listing
	for _, listing := range listings {
		if l.isFavorite(phoneID, listing.ListingID) {
			listing.IsFavorite = true
		}
		result = append(result, listing)
	}
	return result
}

func (l *listingEngine) isFavorite(phoneID string, listingID int) bool {
	fmt.Printf("SELECT favorite_id FROM favorites where phone_id = %s AND listing_id = %d;", phoneID, listingID)
	rows, err := l.sql.Query("SELECT favorite_id FROM favorites where phone_id = $1 AND listing_id = $2;", phoneID, listingID)
	if err != nil {
		return false
	}

	defer rows.Close()

	var favoriteID int
	if rows.Next() {
		err = rows.Scan(&favoriteID)
		if err != nil {
			return false
		}
	}

	if err = rows.Err(); err != nil {
		return false
	}

	if favoriteID != 0 {
		return true
	}

	return false
}
