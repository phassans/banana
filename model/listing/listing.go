package listing

import (
	"bytes"
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/model/business"
	"github.com/phassans/banana/shared"
	"github.com/rs/xlog"
)

type (
	listingEngine struct {
		sql            *sql.DB
		logger         xlog.Logger
		businessEngine business.BusinessEngine
	}

	ListingEngine interface {
		AddListing(listing *shared.Listing) (int, error)
		AddListingImage(businessName string, imagePath string)

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

		GetListingsByBusinessID(businessID int, businessType string) ([]shared.Listing, error)
		GetListingByID(listingID int, businessID int) (shared.Listing, error)
		GetListingInfo(listingID int) (shared.Listing, error)
		GetListingImage() string

		MassageAndPopulateSearchListings([]shared.Listing) ([]shared.SearchListingResult, error)

		DeleteListing(listingId int) error

		ListingEdit(listing *shared.Listing) error
	}
)

func NewListingEngine(psql *sql.DB, logger xlog.Logger, businessEngine business.BusinessEngine) ListingEngine {
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

func (l *listingEngine) AddDietaryRestrictionsToListings(listings []shared.Listing) ([]shared.Listing, error) {
	// get dietary restriction
	var listingsResult []shared.Listing
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

func (l *listingEngine) GetListingImage() string {
	imgRand := random(1, 6)
	return fmt.Sprintf("%s/static/%d.jpg", shared.ImageBaseURL, imgRand)
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
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

/*func (l *listingEngine) GetListingInfo(listingID int) (shared.ListingInfo, error) {
	var listingInfo shared.ListingInfo

	//GetListingByID
	listing, err := l.GetListingByID(listingID, 0)
	if err != nil {
		return shared.ListingInfo{}, err
	}

	if listing.ListingID == 0 {
		return shared.ListingInfo{}, helper.ListingDoesNotExist{ListingID: listingID}
	}

	// add dietary req's
	reqs, err := l.GetDietaryRestriction(listing.ListingID)
	if err != nil {
		return shared.ListingInfo{}, helper.DatabaseError{DBError: err.Error()}
	}
	listing.DietaryRestriction = reqs

	searchListingResult, err := l.MassageAndPopulateSearchListings([]shared.Listing{listing})
	listingInfo.Listing = searchListingResult[0]

	//GetBusinessInfo
	businessInfo, err := l.businessEngine.GetBusinessInfo(listing.BusinessID)
	if err != nil {
		return shared.ListingInfo{}, err
	}
	listingInfo.Business = businessInfo
	return listingInfo, nil
}*/

func (l *listingEngine) GetListingInfo(listingID int) (shared.Listing, error) {
	//var listingInfo shared.Listing

	//GetListingByID
	listing, err := l.GetListingByID(listingID, 0)
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
	listing.DietaryRestriction = reqs

	//GetBusinessInfo
	businessInfo, err := l.businessEngine.GetBusinessInfo(listing.BusinessID)
	if err != nil {
		return shared.Listing{}, err
	}
	listing.Business = businessInfo
	return listing, nil
}

func (l *listingEngine) GetListingByID(listingID int, businessID int) (shared.Listing, error) {

	var whereClause bytes.Buffer
	whereClause.WriteString(fmt.Sprintf(" WHERE listing.listing_id = %d", listingID))
	if businessID != 0 {
		whereClause.WriteString(fmt.Sprintf(" AND business.business_id = %d", businessID))
	}
	query := fmt.Sprintf("%s %s %s;", searchSelect, fromClause, whereClause.String())

	fmt.Println("GetListingByID ", query)

	rows, err := l.sql.Query(query)

	if err != nil {
		return shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var listing shared.Listing
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
			&listing.BusinessName,
			&listing.ListingDate,
		)
		if err != nil {
			return shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
		}
	}

	if err = rows.Err(); err != nil {
		return shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
		return shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
	}

	return listing, nil
}

func (l *listingEngine) GetListingsByBusinessID(businessID int, status string) ([]shared.Listing, error) {
	getListingsQuery := "SELECT listing.title, listing.old_price, listing.new_price, listing.discount, listing.description," +
		"listing.start_date, listing.multiple_days, listing.end_date, listing.start_time, listing.end_time, listing.recurring, " +
		"listing.recurring_end_date, listing.listing_type, " +
		"business.business_id, listing_id, business.name " +
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
	for rows.Next() {
		var listing shared.Listing
		err := rows.Scan(
			&listing.Title,
			&listing.OldPrice,
			&listing.NewPrice,
			&listing.Discount,
			&listing.Description,
			&listing.StartDate,
			&listing.MultipleDays,
			&listing.EndDate,
			&listing.StartTime,
			&listing.EndTime,
			&listing.Recurring,
			&listing.RecurringEndDate,
			&listing.Type,
			&listing.BusinessID,
			&listing.ListingID,
			&listing.BusinessName,
		)
		if err != nil {
			return []shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
		}

		// add dietary req's
		reqs, err := l.GetDietaryRestriction(listing.ListingID)
		if err != nil {
			return []shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
		}
		listing.DietaryRestriction = reqs

		// add recurring listing
		recurring, err := l.GetRecurringListing(listing.ListingID)
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

func (f *listingEngine) isFavorite(phoneID string, listingID int) bool {
	rows, err := f.sql.Query("SELECT favorite_id FROM favorites where phone_id = $1 AND listing_id = $2;", phoneID, listingID)
	if err != nil {
		return false
	}

	defer rows.Close()

	var favoriteID int
	if rows.Next() {
		err := rows.Scan(&favoriteID)
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
