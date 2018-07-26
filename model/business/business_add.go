package business

import (
	"fmt"
	"net/url"

	"github.com/phassans/banana/clients"
	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
	"github.com/rs/xlog"
)

const (
	insertBusinessSQL = "INSERT INTO business(name,phone,website) " +
		"VALUES($1,$2,$3) returning business_id;"

	insertBusinessAddressSQL = "INSERT INTO address(street,city,postal_code,state,country_id,business_id) " +
		"VALUES($1,$2,$3,$4,$5,$6) returning address_id;"

	insertBusinessAddressGEOSQL = "INSERT INTO address_geo(address_id,business_id,latitude,longitude) " +
		"VALUES($1,$2,$3,$4) returning geo_id;"
)

func (b *businessEngine) AddBusiness(
	businessName string,
	phone string,
	website string,
	street string,
	city string,
	postalCode string,
	state string,
	hoursInfo []shared.Hours,
	cuisine []string,
) (int, int, error) {
	// check business name unique
	businessID, err := b.GetBusinessIDFromName(businessName)
	if err != nil {
		return 0, 0, err
	}

	if businessID != 0 {
		return 0, 0, helper.DuplicateEntity{Name: businessName}
	}

	lastInsertBusinessID, err := b.AddBusinessInfo(businessName, phone, website)
	if err != nil {
		return 0, 0, err
	}
	xlog.Infof("AddBusinessInfo success with businessID: %d", lastInsertBusinessID)

	// add business address
	addressID, err := b.AddBusinessAddress(street, city, postalCode, state, lastInsertBusinessID)
	if err != nil {
		return 0, 0, err
	}
	xlog.Infof("AddBusinessAddress success with businessID: %d", lastInsertBusinessID)

	// add business hour
	if err := b.AddBusinessHours(hoursInfo, lastInsertBusinessID); err != nil {
		return 0, 0, err
	}
	xlog.Infof("AddBusinessHours success with businessID: %d", lastInsertBusinessID)

	// add cuisine
	if err := b.AddBusinessCuisine(cuisine, lastInsertBusinessID); err != nil {
		return 0, 0, err
	}
	xlog.Infof("AddBusinessCuisine success with businessID: %d", lastInsertBusinessID)

	b.logger.Infof("business added successfully with id: %d", lastInsertBusinessID)
	return lastInsertBusinessID, addressID, nil
}

func (b *businessEngine) AddBusinessInfo(businessName string, phone string, website string) (int, error) {
	// insert business
	var lastInsertBusinessID int
	err := b.sql.QueryRow(insertBusinessSQL, businessName, phone, website).Scan(&lastInsertBusinessID)
	if err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}
	return lastInsertBusinessID, nil
}

func (b *businessEngine) AddBusinessAddress(
	street string,
	city string,
	postalCode string,
	state string,
	businessID int,
) (int, error) {
	// insert to address table
	var addressID int
	err := b.sql.QueryRow(insertBusinessAddressSQL, street, city, postalCode,
		state, shared.CountryID, businessID).
		Scan(&addressID)
	if err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}
	b.logger.Infof("successfully added address with ID: %d for business: %d", addressID, businessID)

	// add lat, long to database
	geoAddress := fmt.Sprintf("%s,%s,%s", street, city, state)
	err = b.AddGeoInfo(url.QueryEscape(geoAddress), addressID, businessID)
	if err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}

	return addressID, nil
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

func (l *businessEngine) AddBusinessHours(days []shared.Hours, businessID int) error {
	business, err := l.GetBusinessFromID(businessID)
	if err != nil {
		return err
	}

	if business.Name == "" {
		return helper.BusinessError{Message: fmt.Sprintf("business with id %d does not exist", businessID)}
	}

	for _, day := range days {
		if day.Day != "" {
			if day.OpenTimeSessionOne != "" && day.CloseTimeSessionOne != "" {
				if err := l.addHours(day.Day, day.OpenTimeSessionOne, day.CloseTimeSessionOne, businessID); err != nil {
					return err
				}
			}

			if day.OpenTimeSessionTwo != "" && day.CloseTimeSessionTwo != "" {
				if err := l.addHours(day.Day, day.OpenTimeSessionTwo, day.CloseTimeSessionTwo, businessID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (l *businessEngine) addHours(day string, openTime string, closeTime string, businessID int) error {
	addBusinessHoursSQL := "INSERT INTO business_hours(business_id,day,open_time,close_time) " +
		"VALUES($1,$2,$3,$4);"

	_, err := l.sql.Query(addBusinessHoursSQL, businessID, day, openTime, closeTime)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	l.logger.Infof("add hours succesfull for businessID:%d", businessID)
	return nil
}

func (l *businessEngine) AddBusinessCuisine(cuisines []string, businessID int) error {
	business, err := l.GetBusinessFromID(businessID)
	if err != nil {
		return err
	}

	if business.Name == "" {
		return helper.BusinessError{Message: fmt.Sprintf("business with id %d does not exist", businessID)}
	}

	for _, cuisine := range cuisines {
		if err := l.addCuisine(businessID, cuisine); err != nil {
			return err
		}
	}

	return nil
}

func (l *businessEngine) addCuisine(businessID int, cuisine string) error {
	addBusinessCuisineSQL := "INSERT INTO business_cuisine(business_id,cuisine) " +
		"VALUES($1,$2);"

	_, err := l.sql.Query(addBusinessCuisineSQL, businessID, cuisine)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	l.logger.Infof("add cuisine successful for businessID:%d", businessID)
	return nil
}
