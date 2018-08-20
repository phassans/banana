package favourite

import (
	"bytes"
	"database/sql"
	"fmt"
	"time"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/model/business"
	"github.com/phassans/banana/model/common"
	"github.com/phassans/banana/model/listing"
	"github.com/phassans/banana/shared"
	"github.com/rs/zerolog"
)

type (
	favoriteEngine struct {
		sql            *sql.DB
		logger         zerolog.Logger
		businessEngine business.BusinessEngine
		listingEngine  listing.ListingEngine
	}

	// FavoriteEngine interface which holds all methods
	FavoriteEngine interface {
		AddFavorite(phoneID string, listingID int, listingDateID int) error
		DeleteFavorite(phoneID string, listingID int, listingDateID int) error
		GetAllFavorites(phoneID string, sortBy string, latitude float64, longitude float64) ([]shared.SearchListingResult, error)
	}
)

// NewFavoriteEngine returns an instance of favoriteEngine
func NewFavoriteEngine(psql *sql.DB, logger zerolog.Logger, businessEngine business.BusinessEngine, listingEngine listing.ListingEngine) FavoriteEngine {
	return &favoriteEngine{psql, logger, businessEngine, listingEngine}
}

func (f *favoriteEngine) AddFavorite(phoneID string, listingID int, listingDateID int) error {
	listing, err := f.listingEngine.GetListingByID(listingID, 0, listingDateID)
	if err != nil {
		return err
	}

	if listing.ListingID == 0 {
		return helper.ListingDoesNotExist{ListingID: listingID}
	}

	var favoriteID int
	err = f.sql.QueryRow("INSERT INTO favorites(phone_id,listing_id,listing_date_id,favorite_add_date) "+
		"VALUES($1,$2,$3,$4) returning favorite_id;",
		phoneID, listingID, listingDateID, time.Now()).Scan(&favoriteID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	f.logger.Info().Msgf("successfully added a favorites with ID: %d", favoriteID)
	return nil
}

func (f *favoriteEngine) DeleteFavorite(phoneID string, listingID int, listingDateID int) error {
	sqlStatement := `DELETE FROM favorites WHERE phone_id = $1 AND listing_id = $2 AND listing_date_id=$3;`
	//f.logger.Info().Msgf("deleting favorites with query: %s and listing: %d", sqlStatement, listingID)

	_, err := f.sql.Exec(sqlStatement, phoneID, listingID, listingDateID)
	return err
}

/*func (f *favoriteEngine) GetAllFavorites(phoneID string, sortBy string, latitude float64, longitude float64) ([]shared.SearchListingResult, error) {
	listings, err := f.GetListingsPhoneID(phoneID)
	if err != nil {
		return nil, err
	}

	sortEngine := listing.NewSortListingEngine(listings, sortBy, shared.CurrentLocation{Latitude: latitude, Longitude: longitude}, f.sql)
	listings, err = sortEngine.SortListings()
	if err != nil {
		return nil, err
	}
	f.logger.Info().Msgf("done sorting the listings in favorite. listings count: %d", len(listings))

	return f.listingEngine.MassageAndPopulateSearchListings(listings)
}*/

func (f *favoriteEngine) GetAllFavorites(phoneID string, sortBy string, latitude float64, longitude float64) ([]shared.SearchListingResult, error) {
	/*favorites, err := f.GetAllFavoritesIDs(phoneID)
	if err != nil {
		return nil, err
	}

	var listings []shared.Listing
	for _, favorite := range favorites {
		listing, err := f.listingEngine.GetListingByID(favorite.ListingID, 0, favorite.ListingDateID)
		if err != nil {
			return nil, err
		}
		listing.ListingImage = listing.ImageLink
		listing.Favorite = &favorite
		listings = append(listings, listing)
	}*/

	listings, err := f.GetListingsPhoneID(phoneID)
	if err != nil {
		return nil, err
	}

	sortEngine := listing.NewSortListingEngine(listings, sortBy, shared.CurrentLocation{Latitude: latitude, Longitude: longitude}, f.sql)
	listings, err = sortEngine.SortListings(false)
	if err != nil {
		return nil, err
	}
	f.logger.Info().Msgf("done sorting the listings in favorite. listings count: %d", len(listings))

	return f.listingEngine.MassageAndPopulateSearchListings(listings)

}

func (f *favoriteEngine) GetAllFavoritesIDs(phoneID string) ([]shared.Favorite, error) {
	rows, err := f.sql.Query("SELECT favorite_id,listing_id,listing_date_id,favorite_add_date FROM favorites where phone_id = $1;", phoneID)
	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var favorites []shared.Favorite
	for rows.Next() {
		var favorite shared.Favorite
		err = rows.Scan(&favorite.FavoriteID, &favorite.ListingID, &favorite.ListingDateID, &favorite.FavoriteAddDate)
		if err != nil {
			return nil, helper.DatabaseError{DBError: err.Error()}
		}
		favorites = append(favorites, favorite)
	}

	if err = rows.Err(); err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	return favorites, nil
}

const (
/*searchSelect = "SELECT listing.title as title, listing.old_price as old_price, listing.new_price as new_price," +
"listing.discount as discount, listing.discount_description as discount_description, listing.description as description, listing.start_date as start_date," +
"listing.end_date as end_date, listing.start_time as start_time, listing.end_time as end_time," +
"listing.multiple_days as multiple_days," +
"listing.recurring as recurring, listing.recurring_end_date as recurring_date, listing.listing_type as listing_type, " +
"listing.business_id as business_id, listing.listing_id as listing_id, " +
"business.name as bname, " +
"listing_date.listing_date_id as listing_date_id, listing_date.listing_date as listing_date, " +
"listing_image.path as path, " +
"favorites.favorite_id as favorite_id, favorites.favorite_add_date as favorite_add_date "*/

)

func (f *favoriteEngine) GetListingsPhoneID(phoneID string) ([]shared.Listing, error) {

	var whereClause bytes.Buffer
	whereClause.WriteString(fmt.Sprintf(" WHERE favorites.phone_id = '%s'", phoneID))

	selectFields := fmt.Sprintf("%s, %s, %s, %s, %s", common.ListingFields, common.ListingBusinessFields, common.ListingDateFields, common.ListingImageFields, common.FavoriteFields)
	fmt.Println("selectFields: ", selectFields)

	query := fmt.Sprintf("%s %s %s", selectFields, common.FromClauseFavorites, whereClause.String())
	rows, err := f.sql.Query(query)

	fmt.Println("GetListingsPhoneID ", query)

	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var listings []shared.Listing
	for rows.Next() {
		var sqlEndDate sql.NullString
		var sqlRecurringEndDate sql.NullString
		var sqlFavoriteAddDate sql.NullString
		var listing shared.Listing
		var fid int
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
			&fid,
			&sqlFavoriteAddDate,
		)

		listing.Favorite = &shared.Favorite{fid, listing.ListingID, listing.ListingDateID, sqlFavoriteAddDate.String}
		listing.StartDate = sqlEndDate.String
		listing.RecurringEndDate = sqlRecurringEndDate.String

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
