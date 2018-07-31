package favourite

import (
	"database/sql"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/model/business"
	"github.com/phassans/banana/model/listing"
	"github.com/phassans/banana/shared"
	"github.com/rs/xlog"
)

type (
	favoriteEngine struct {
		sql            *sql.DB
		logger         xlog.Logger
		businessEngine business.BusinessEngine
		listingEngine  listing.ListingEngine
	}

	FavoriteEngine interface {
		AddFavorite(phoneID string, listingID int) error
		DeleteFavorite(phoneID string, listingID int) error
		GetAllFavorites(phoneID string) ([]shared.SearchListingResult, error)
	}
)

func NewFavoriteEngine(psql *sql.DB, logger xlog.Logger, businessEngine business.BusinessEngine, listingEngine listing.ListingEngine) FavoriteEngine {
	return &favoriteEngine{psql, logger, businessEngine, listingEngine}
}

func (f *favoriteEngine) AddFavorite(phoneID string, listingID int) error {
	listing, err := f.listingEngine.GetListingByID(listingID, 0)
	if err != nil {
		return err
	}

	if listing.ListingID == 0 {
		return helper.ListingDoesNotExist{ListingID: listingID}
	}

	var favoriteID int
	err = f.sql.QueryRow("INSERT INTO favorites(phone_id,listing_id) "+
		"VALUES($1,$2) returning favorite_id;",
		phoneID, listingID).Scan(&favoriteID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	f.logger.Infof("successfully added a favorites with ID: %d", favoriteID)
	return nil
}

func (f *favoriteEngine) DeleteFavorite(phoneID string, listingID int) error {
	sqlStatement := `DELETE FROM favorites WHERE phone_id = $1 AND listing_id = $2;`
	f.logger.Infof("deleting favorites with query: %s and listing: %d", sqlStatement, listingID)

	_, err := f.sql.Exec(sqlStatement, phoneID, listingID)
	return err
}

func (f *favoriteEngine) GetAllFavorites(phoneID string) ([]shared.SearchListingResult, error) {
	IDs, err := f.GetAllFavoritesIDs(phoneID)
	if err != nil {
		return nil, err
	}

	var listings []shared.Listing
	for _, id := range IDs {
		listing, err := f.listingEngine.GetListingByID(id, 0)
		if err != nil {
			return nil, err
		}

		imageLink, err := f.listingEngine.GetListingImage(id)
		if err != nil {
			return nil, err
		}
		listing.ListingImage = imageLink
		listings = append(listings, listing)
	}

	return f.listingEngine.MassageAndPopulateSearchListings(listings)
}

func (f *favoriteEngine) GetAllFavoritesIDs(phoneID string) ([]int, error) {
	rows, err := f.sql.Query("SELECT listing_id FROM favorites where phone_id = $1;", phoneID)
	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var listingIDs []int
	for rows.Next() {
		var listingID int
		err := rows.Scan(&listingID)
		if err != nil {
			return nil, helper.DatabaseError{DBError: err.Error()}
		}
		listingIDs = append(listingIDs, listingID)
	}

	if err = rows.Err(); err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	return listingIDs, nil
}
