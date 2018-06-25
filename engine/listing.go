package engine

import (
	"database/sql"
	"time"

	"github.com/pshassans/banana/helper"
	"github.com/rs/xlog"
)

type listingEngine struct {
	sql    *sql.DB
	logger xlog.Logger
}

type ListingEngine interface {
	AddBusiness(businessName string, phone string, website string) (int, error)
	GetBusinessIDFromName(businessName string) (int, error)
	AddOwner(firstName string, lastName string, phone string, email string, businessName string) error
	AddBusinessAddress(line1 string, line2 string, city string, postalCode string, state string, country string, businessName string, otherDetails string) error
	AddListing(title string, description string, price float64, startTime string, endTime string, businessName string) error
}

func NewListingEngine(psql *sql.DB, logger xlog.Logger) ListingEngine {
	return &listingEngine{psql, logger}
}

func (l *listingEngine) AddBusiness(businessName string, phone string, website string) (int, error) {
	businessID, err := l.GetBusinessIDFromName(businessName)
	if err != nil {
		return 0, err
	}

	if businessID != 0 {
		return 0, helper.DuplicateEntity{BusinessName: businessName}
	}

	var lastInsertBusinessID int

	err = l.sql.QueryRow("INSERT INTO business(name,phone,website) "+
		"VALUES($1,$2,$3) returning business_id;",
		businessName, phone, website).Scan(&lastInsertBusinessID)
	if err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}

	l.logger.Infof("last inserted id: %d", lastInsertBusinessID)

	return lastInsertBusinessID, nil
}

func (l *listingEngine) GetBusinessIDFromName(businessName string) (int, error) {
	rows, err := l.sql.Query("SELECT business_id FROM business where name = $1;", businessName)
	if err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var id int
	if rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return 0, helper.DatabaseError{DBError: err.Error()}
		}
	}

	if err = rows.Err(); err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}

	return id, nil
}

func (l *listingEngine) AddOwner(firstName string, lastName string, phone string, email string, businessName string) error {
	businessID, err := l.GetBusinessIDFromName(businessName)
	if err != nil {
		return err
	}

	if businessID == 0 {
		return helper.BusinessDoesNotExist{BusinessName: businessName}
	}

	var ownerID int

	err = l.sql.QueryRow("INSERT INTO owner(first_name,last_name,phone,email,business_id) "+
		"VALUES($1,$2,$3,$4,$5) returning owner_id;",
		firstName, lastName, phone, email, businessID).Scan(&ownerID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	l.logger.Infof("successfully added a user with ID: %d", ownerID)

	return nil
}

func (l *listingEngine) AddBusinessAddress(line1 string, line2 string, city string, postalCode string, state string, country string, businessName string, otherDetails string) error {
	businessID, err := l.GetBusinessIDFromName(businessName)
	if err != nil {
		return err
	}

	if businessID == 0 {
		return helper.BusinessDoesNotExist{BusinessName: businessName}
	}

	var addressID int
	countryID := 1

	err = l.sql.QueryRow("INSERT INTO address(line1,line2,city,postal_code,state,country_id,business_id,other_details) "+
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8) returning address_id;",
		line1, line2, city, postalCode, state, countryID, businessID, otherDetails).Scan(&addressID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	l.logger.Infof("successfully added a address with ID: %d for business: %s", addressID, businessName)

	return nil
}

func (l *listingEngine) AddListing(title string, description string, price float64, startTime string, endTime string, businessName string) error {
	businessID, err := l.GetBusinessIDFromName(businessName)
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
