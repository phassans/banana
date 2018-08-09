package favourite

import (
	"database/sql"

	"time"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/model/business"
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
		AddFavorite(phoneID string, listingID int) error
		DeleteFavorite(phoneID string, listingID int) error
		GetAllFavorites(phoneID string, sortBy string, latitude float64, longitude float64) ([]shared.SearchListingResult, error)
	}
)

// NewFavoriteEngine returns an instance of favoriteEngine
func NewFavoriteEngine(psql *sql.DB, logger zerolog.Logger, businessEngine business.BusinessEngine, listingEngine listing.ListingEngine) FavoriteEngine {
	return &favoriteEngine{psql, logger, businessEngine, listingEngine}
}

func (f *favoriteEngine) AddFavorite(phoneID string, listingID int) error {
	listing, err := f.listingEngine.GetListingByID(listingID, 0, 0)
	if err != nil {
		return err
	}

	if listing.ListingID == 0 {
		return helper.ListingDoesNotExist{ListingID: listingID}
	}

	var favoriteID int
	err = f.sql.QueryRow("INSERT INTO favorites(phone_id,listing_id,favorite_add_date) "+
		"VALUES($1,$2,$3) returning favorite_id;",
		phoneID, listingID, time.Now()).Scan(&favoriteID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	f.logger.Info().Msgf("successfully added a favorites with ID: %d", favoriteID)
	return nil
}

func (f *favoriteEngine) DeleteFavorite(phoneID string, listingID int) error {
	sqlStatement := `DELETE FROM favorites WHERE phone_id = $1 AND listing_id = $2;`
	f.logger.Info().Msgf("deleting favorites with query: %s and listing: %d", sqlStatement, listingID)

	_, err := f.sql.Exec(sqlStatement, phoneID, listingID)
	return err
}

func (f *favoriteEngine) GetAllFavorites(phoneID string, sortBy string, latitude float64, longitude float64) ([]shared.SearchListingResult, error) {
	favorites, err := f.GetAllFavoritesIDs(phoneID)
	if err != nil {
		return nil, err
	}

	var listings []shared.Listing
	for _, favorite := range favorites {
		listing, err := f.listingEngine.GetListingByID(favorite.ListingID, 0, 0)
		if err != nil {
			return nil, err
		}
		listing.Favorite = &favorite

		imageLink, err := f.listingEngine.GetListingImage(favorite.ListingID)
		if err != nil {
			return nil, err
		}
		listing.ListingImage = imageLink
		listings = append(listings, listing)
	}

	sortEngine := listing.NewSortListingEngine(listings, sortBy, shared.CurrentLocation{Latitude: latitude, Longitude: longitude}, f.sql)
	listings, err = sortEngine.SortListings()
	if err != nil {
		return nil, err
	}
	f.logger.Info().Msgf("done sorting the listings in favorite. listings count: %d", len(listings))

	return f.listingEngine.MassageAndPopulateSearchListings(listings)
}

func (f *favoriteEngine) GetAllFavoritesIDs(phoneID string) ([]shared.Favorite, error) {
	rows, err := f.sql.Query("SELECT favorite_id,listing_id,favorite_add_date FROM favorites where phone_id = $1;", phoneID)
	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var favorites []shared.Favorite
	for rows.Next() {
		var favorite shared.Favorite
		err = rows.Scan(&favorite.FavoriteID, &favorite.ListingID, &favorite.FavoriteAddDate)
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
