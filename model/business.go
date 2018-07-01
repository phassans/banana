package model

import (
	"database/sql"

	"github.com/phassans/banana/clients"
	"github.com/phassans/banana/helper"
	"github.com/rs/xlog"
)

type businessEngine struct {
	sql    *sql.DB
	logger xlog.Logger
}

const (
	insertBusinessSQL = "INSERT INTO business(name,phone,website) " +
		"VALUES($1,$2,$3) returning business_id;"

	insertBusinessAddressSQL = "INSERT INTO address(line1,line2,city,postal_code,state,country_id,business_id,other_details) " +
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8) returning address_id;"

	insertBusinessAddressGEOSQL = "INSERT INTO address_geo(address_id,business_id,latitude,longitude) " +
		"VALUES($1,$2,$3,$4) returning geo_id;"
)

type BusinessEngine interface {
	AddBusiness(
		businessName string,
		phone string,
		website string,
		street string,
		city string,
		postalCode string,
		state string,
		country string,
	) (int, error)

	AddBusinessAddress(
		street string,
		city string,
		postalCode string,
		state string,
		country string,
		businessID int,
	) error
	AddGeoInfo(address string, addressID int, businessID int) error
	GetBusinessIDFromName(businessName string) (int, error)

	// TBD
	AddBusinessHours(businessName string, day string, openTime string, closeTime string)
	AddBusinessImage(businessName string, imagePath string)
}

func NewBusinessEngine(psql *sql.DB, logger xlog.Logger) BusinessEngine {
	return &businessEngine{psql, logger}
}

func (l *businessEngine) AddBusiness(
	businessName string,
	phone string,
	website string,
	street string,
	city string,
	postalCode string,
	state string,
	country string,
) (int, error) {
	businessID, err := l.GetBusinessIDFromName(businessName)
	if err != nil {
		return 0, err
	}

	if businessID != 0 {
		return 0, helper.DuplicateEntity{BusinessName: businessName}
	}

	var lastInsertBusinessID int

	err = l.sql.QueryRow(insertBusinessSQL, businessName, phone, website).
		Scan(&lastInsertBusinessID)
	if err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}

	l.logger.Infof("last inserted id: %d", lastInsertBusinessID)

	return lastInsertBusinessID, nil
}

func (l *businessEngine) AddBusinessHours(businessName string, day string, openTime string, closeTime string) {
	return
}

func (l *businessEngine) AddBusinessImage(businessName string, imagePath string) {
	return
}

func (l *businessEngine) GetBusinessIDFromName(businessName string) (int, error) {
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

func (l *businessEngine) AddBusinessAddress(
	street string,
	city string,
	postalCode string,
	state string,
	country string,
	businessID int,
) error {
	/*businessID, err := l.GetBusinessIDFromName(businessName)
	if err != nil {
		return err
	}

	if businessID == 0 {
		return helper.BusinessDoesNotExist{BusinessName: businessName}
	}

	var addressID int
	countryID := 1

	err = l.sql.QueryRow(insertBusinessAddressSQL, line1, line2, city, postalCode,
		state, countryID, businessID, otherDetails).
		Scan(&addressID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	l.logger.Infof("successfully added a address with ID: %d for business: %s", addressID, businessName)

	geoAddress := fmt.Sprintf("%s,%s,%s,%s", line1, line2, city, state)

	// add lat, long to database
	err = l.AddGeoInfo(url.QueryEscape(geoAddress), addressID, businessID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}*/

	return nil
}

func (l *businessEngine) AddGeoInfo(address string, addressID int, businessID int) error {
	resp, err := clients.GetLatLong(address)

	var geoID int

	err = l.sql.QueryRow(insertBusinessAddressGEOSQL, addressID, businessID, resp.Lat, resp.Lon).
		Scan(&geoID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	l.logger.Infof("successfully added a geoLocation %s for address: %s", geoID, address)

	return err
}
