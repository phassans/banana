package engine

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

func NewListingEngine(psql *sql.DB, logger xlog.Logger, businessEngine BusinessEngine) ListingEngine {
	return &listingEngine{psql, logger, businessEngine}
}

func (l *listingEngine) AddListing(title string, description string, price float64, startTime string, endTime string, businessName string) error {
	businessID, err := l.businessEngine.GetBusinessIDFromName(businessName)
	if err != nil {
		return err
	}

	if businessID == 0 {
		return helper.BusinessDoesNotExist{BusinessName: businessName}
	}

	var listingID int

	err = l.sql.QueryRow("INSERT INTO listing(title,description,price,start_time,end_time,business_id) "+
		"VALUES($1,$2,$3,$4,$5,$6) returning listing_id;",
		title, description, price, time.Now(), time.Now(), businessID).Scan(&listingID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	l.logger.Infof("successfully added a listing %s for business: %s", title, businessName)

	return nil
}
