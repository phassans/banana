package model

import (
	"database/sql"

	"github.com/phassans/banana/helper"
	"github.com/rs/xlog"
)

type (
	favoriteEngine struct {
		sql            *sql.DB
		logger         xlog.Logger
		businessEngine BusinessEngine
		listingEngine  ListingEngine
	}

	FavoriteEngine interface {
		AddFavorite(phoneID string, listingID int) error
		DeleteFavorite(phoneID string, listingID int) error
		GetAllFavorites(phoneID string) ([]Listing, error)
	}
)

func NewFavoriteEngine(psql *sql.DB, logger xlog.Logger, businessEngine BusinessEngine, listingEngine ListingEngine) FavoriteEngine {
	return &favoriteEngine{psql, logger, businessEngine, listingEngine}
}

func (f *favoriteEngine) AddFavorite(phoneID string, listingID int) error {
	var favoriteID int
	err := f.sql.QueryRow("INSERT INTO favorites(phone_id,listing_id) "+
		"VALUES($1,$2) returning favorite_id;",
		phoneID, listingID).Scan(&favoriteID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	f.logger.Infof("successfully added a user with ID: %d", favoriteID)
	return nil
}

func (f *favoriteEngine) DeleteFavorite(phoneID string, listingID int) error {
	sqlStatement := `DELETE FROM favorites WHERE phone_id = $1 AND listing_id = $2;`
	f.logger.Infof("deleting favorites with query: %s and listing: %d", sqlStatement, listingID)

	_, err := f.sql.Exec(sqlStatement, phoneID, listingID)
	return err
}

func (f *favoriteEngine) GetAllFavorites(phoneID string) ([]Listing, error) {
	IDs, err := f.GetAllFavoritesIDs(phoneID)
	if err != nil {
		return nil, err
	}

	var listings []Listing
	for _, id := range IDs {
		listing, err := f.listingEngine.GetListingByID(id)
		if err != nil {
			return nil, err
		}
		listings = append(listings, listing)
	}

	return listings, nil
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
