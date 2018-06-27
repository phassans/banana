package model

import (
	"database/sql"
	"time"

	"github.com/pshassans/banana/helper"
	"github.com/rs/xlog"
)

type listingEngine struct {
	sql            *sql.DB
	logger         xlog.Logger
	businessEngine BusinessEngine
}

const insertListingSQL = "INSERT INTO listing(title, description, old_price, new_price, " +
	"listing_date, start_time, end_time, business_id, recurring, listing_create_date) " +
	"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning listing_id"

type ListingEngine interface {
	AddListing(
		title string,
		description string,
		oldPrice float64,
		newPrice float64,
		listingDate string,
		startTime string,
		endTime string,
		recurring bool,
		businessName string,
	) error
	AddListingImage(businessName string, imagePath string)
	AddListingDietaryRestrictions(listingTitle string, dietaryRestriction string)
	AddListingRecurringInfo(listingTitle string, day string, startTime string, endTime string)
}

func NewListingEngine(psql *sql.DB, logger xlog.Logger, businessEngine BusinessEngine) ListingEngine {
	return &listingEngine{psql, logger, businessEngine}
}

func (l *listingEngine) AddListing(
	title string,
	description string,
	oldPrice float64,
	newPrice float64,
	listingDate string,
	startTime string,
	endTime string,
	recurring bool,
	businessName string,
) error {
	businessID, err := l.businessEngine.GetBusinessIDFromName(businessName)
	if err != nil {
		return err
	}

	if businessID == 0 {
		return helper.BusinessDoesNotExist{BusinessName: businessName}
	}

	var listingID int

	err = l.sql.QueryRow(insertListingSQL, title, description, oldPrice, newPrice,
		listingDate, time.Now(), time.Now(), businessID, recurring, time.Now()).
		Scan(&listingID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	l.logger.Infof("successfully added a listing %s for business: %s", title, businessName)

	return nil
}

func (l *listingEngine) AddListingImage(businessName string, imagePath string) {
	return
}

func (l *listingEngine) AddListingDietaryRestrictions(listingTitle string, dietaryRestriction string) {
	return
}

func (l *listingEngine) AddListingRecurringInfo(listingTitle string, day string, startTime string, endTime string) {
	return
}
