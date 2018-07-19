package business

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/phassans/banana/clients"
	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
	"github.com/rs/xlog"
)

type businessEngine struct {
	sql    *sql.DB
	logger xlog.Logger
}

const (
	insertBusinessSQL = "INSERT INTO business(name,phone,website) " +
		"VALUES($1,$2,$3) returning business_id;"

	insertBusinessAddressSQL = "INSERT INTO address(street,city,postal_code,state,country_id,business_id) " +
		"VALUES($1,$2,$3,$4,$5,$6) returning address_id;"

	insertBusinessAddressGEOSQL = "INSERT INTO address_geo(address_id,business_id,latitude,longitude) " +
		"VALUES($1,$2,$3,$4) returning geo_id;"
)

type BusinessEngine interface {
	// Add
	AddBusiness(
		businessName string,
		phone string,
		website string,
		street string,
		city string,
		postalCode string,
		state string,
		country string,
		hoursInfo []shared.Hours,
		cuisine []string,
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
	AddBusinessHours([]shared.Hours, int) error
	AddBusinessCuisine(cuisines []string, businessID int) error

	// Select
	GetBusinessIDFromName(businessName string) (int, error)
	GetBusinessFromID(businessID int) (shared.Business, error)
	GetBusinessAddressFromID(businessID int) (shared.BusinessAddress, error)
	GetBusinessInfo(businessID int) (shared.BusinessInfo, error)
	GetAllBusiness() ([]shared.Business, error)

	// Delete
	DeleteBusinessFromID(businessID int) error
	DeleteBusinessAddressFromID(businessID int) error
}

func NewBusinessEngine(psql *sql.DB, logger xlog.Logger) BusinessEngine {
	return &businessEngine{psql, logger}
}

func (b *businessEngine) GetAllBusiness() ([]shared.Business, error) {
	rows, err := b.sql.Query("SELECT business_id, name, phone, website FROM business;")
	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var allBusiness []shared.Business
	for rows.Next() {
		var business shared.Business
		err := rows.Scan(&business.BusinessID, &business.Name, &business.Phone, &business.Website)
		if err != nil {
			return nil, helper.DatabaseError{DBError: err.Error()}
		}
		allBusiness = append(allBusiness, business)
	}

	if err = rows.Err(); err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	return allBusiness, nil
}

func (b *businessEngine) AddBusiness(
	businessName string,
	phone string,
	website string,
	street string,
	city string,
	postalCode string,
	state string,
	country string,
	hoursInfo []shared.Hours,
	cuisine []string,
) (int, error) {
	// check business name unique
	businessID, err := b.GetBusinessIDFromName(businessName)
	if err != nil {
		return 0, err
	}

	if businessID != 0 {
		return 0, helper.DuplicateEntity{Name: businessName}
	}

	// insert business
	var lastInsertBusinessID int
	err = b.sql.QueryRow(insertBusinessSQL, businessName, phone, website).Scan(&lastInsertBusinessID)
	if err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}

	// add business address
	if err = b.AddBusinessAddress(street, city, postalCode, state, country, lastInsertBusinessID); err != nil {
		return 0, err
	}

	// add business hour
	if err := b.AddBusinessHours(hoursInfo, lastInsertBusinessID); err != nil {
		return 0, err
	}

	// add cuisine
	if err := b.AddBusinessCuisine(cuisine, lastInsertBusinessID); err != nil {
		return 0, err
	}

	b.logger.Infof("business addded successfully with id: %d", lastInsertBusinessID)

	return lastInsertBusinessID, nil
}

func (b *businessEngine) AddBusinessAddress(
	street string,
	city string,
	postalCode string,
	state string,
	country string,
	businessID int,
) error {
	// insert to address table
	var addressID int
	err := b.sql.QueryRow(insertBusinessAddressSQL, street, city, postalCode,
		state, shared.CountryID, businessID).
		Scan(&addressID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}
	b.logger.Infof("successfully added address with ID: %d for business: %d", addressID, businessID)

	// add lat, long to database
	geoAddress := fmt.Sprintf("%s,%s,%s", street, city, state)
	err = b.AddGeoInfo(url.QueryEscape(geoAddress), addressID, businessID)
	if err != nil {
		// cleanup
		go func() {
			b.DeleteBusinessAddressFromID(addressID)
			b.DeleteBusinessFromID(businessID)
		}()
		return helper.DatabaseError{DBError: err.Error()}
	}

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
				if err := l.AddHours(day.Day, day.OpenTimeSessionOne, day.CloseTimeSessionOne, businessID); err != nil {
					return err
				}
			}

			if day.OpenTimeSessionTwo != "" && day.CloseTimeSessionTwo != "" {
				if err := l.AddHours(day.Day, day.OpenTimeSessionTwo, day.CloseTimeSessionTwo, businessID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (l *businessEngine) AddHours(day string, openTime string, closeTime string, businessID int) error {
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
		if err := l.AddCuisine(businessID, cuisine); err != nil {
			return err
		}
	}

	return nil
}

func (l *businessEngine) AddCuisine(businessID int, cuisine string) error {
	addBusinessCuisineSQL := "INSERT INTO business_cuisine(business_id,cuisine) " +
		"VALUES($1,$2);"

	_, err := l.sql.Query(addBusinessCuisineSQL, businessID, cuisine)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	l.logger.Infof("add cuisine successful for businessID:%d", businessID)
	return nil
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

func (l *businessEngine) GetBusinessFromID(businessID int) (shared.Business, error) {
	rows, err := l.sql.Query("SELECT business_id, name, phone, website FROM business where business_id = $1;", businessID)
	if err != nil {
		return shared.Business{}, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var business shared.Business
	if rows.Next() {
		err := rows.Scan(&business.BusinessID, &business.Name, &business.Phone, &business.Website)
		if err != nil {
			return shared.Business{}, helper.DatabaseError{DBError: err.Error()}
		}
	}

	if err = rows.Err(); err != nil {
		return shared.Business{}, helper.DatabaseError{DBError: err.Error()}
	}

	return business, nil
}

func (l *businessEngine) GetBusinessAddressFromID(businessID int) (shared.BusinessAddress, error) {
	rows, err := l.sql.Query("SELECT street, city, postal_code, state FROM address where business_id = $1;", businessID)
	if err != nil {
		return shared.BusinessAddress{}, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var businessAddress shared.BusinessAddress
	if rows.Next() {
		err := rows.Scan(&businessAddress.Street, &businessAddress.City, &businessAddress.PostalCode, &businessAddress.State)
		if err != nil {
			return shared.BusinessAddress{}, helper.DatabaseError{DBError: err.Error()}
		}
	}

	if err = rows.Err(); err != nil {
		return shared.BusinessAddress{}, helper.DatabaseError{DBError: err.Error()}
	}

	return businessAddress, nil
}

func (l *businessEngine) GetBusinessHoursFromID(businessID int) ([]shared.Bhour, error) {
	rows, err := l.sql.Query("SELECT day, open_time, close_time FROM business_hours where business_id = $1;", businessID)
	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var businessHours []shared.Bhour
	for rows.Next() {
		var businessHour shared.Bhour
		err := rows.Scan(&businessHour.Day, &businessHour.OpenTime, &businessHour.CloseTime)
		if err != nil {
			return nil, helper.DatabaseError{DBError: err.Error()}
		}
		businessHours = append(businessHours, businessHour)
	}

	if err = rows.Err(); err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	return businessHours, nil
}

func (l *businessEngine) GetBusinessInfo(businessID int) (shared.BusinessInfo, error) {
	business, err := l.GetBusinessFromID(businessID)
	if err != nil {
		return shared.BusinessInfo{}, err
	}

	businessAddress, err := l.GetBusinessAddressFromID(businessID)
	if err != nil {
		return shared.BusinessInfo{}, err
	}

	businessHours, err := l.GetBusinessHoursFromID(businessID)
	if err != nil {
		return shared.BusinessInfo{}, err
	}

	return shared.BusinessInfo{Business: business, BusinessAddress: businessAddress, Hours: businessHours}, nil

}

func (b *businessEngine) DeleteBusinessAddressFromID(addressID int) error {
	sqlStatement := `DELETE FROM address WHERE address_id = $1;`
	b.logger.Infof("deleting address with query: %s and business_id: %d", sqlStatement, addressID)

	_, err := b.sql.Exec(sqlStatement, addressID)
	return err
}

func (b *businessEngine) DeleteBusinessFromID(businessID int) error {
	sqlStatement := `DELETE FROM business WHERE business_id = $1;`
	b.logger.Infof("deleting business with query: %s and business_id: %d", sqlStatement, businessID)

	_, err := b.sql.Exec(sqlStatement, businessID)
	return err
}
