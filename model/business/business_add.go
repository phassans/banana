package business

import (
	"fmt"

	"github.com/phassans/banana/clients"
	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
	"github.com/rs/xlog"
)

const (
	insertBusinessSQL = "INSERT INTO business(name,phone,website) " +
		"VALUES($1,$2,$3) returning business_id;"

	insertBusinessAddressSQL = "INSERT INTO business_address(business_id,street,city,postal_code,state,country_id,latitude,longitude) " +
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8) returning address_id;"
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
	userID int,
) (int, int, error) {

	if userID != 0 {
		bUser, err := b.userEngine.UserGet(userID)
		if err != nil {
			return 0, 0, err
		}

		if bUser.UserID == 0 {
			return 0, 0, helper.UserError{Message: fmt.Sprintf("user with ID:%d does not exist", bUser.UserID)}
		}
	}

	// check business name unique
	businessID, err := b.GetBusinessIDFromName(businessName)
	if err != nil {
		return 0, 0, err
	}

	if businessID != 0 {
		return 0, 0, helper.DuplicateEntity{Name: businessName}
	}

	lastInsertBusinessID, err := b.addBusinessInfo(businessName, phone, website)
	if err != nil {
		return 0, 0, err
	}
	xlog.Infof("addBusinessInfo success with businessID: %d", lastInsertBusinessID)

	// add business address
	addressID, err := b.addBusinessAddress(street, city, postalCode, state, lastInsertBusinessID)
	if err != nil {
		return 0, 0, err
	}
	xlog.Infof("addBusinessAddress success with businessID: %d", lastInsertBusinessID)

	// add business hour
	if err := b.addBusinessHours(hoursInfo, lastInsertBusinessID); err != nil {
		return 0, 0, err
	}
	xlog.Infof("addBusinessHours success with businessID: %d", lastInsertBusinessID)

	// add cuisine
	if err := b.addBusinessCuisine(cuisine, lastInsertBusinessID); err != nil {
		return 0, 0, err
	}
	xlog.Infof("addBusinessCuisine success with businessID: %d", lastInsertBusinessID)

	// associate user to business
	if userID != 0 {
		if err := b.associateUserToBusiness(lastInsertBusinessID, userID); err != nil {
			return 0, 0, err
		}
		xlog.Infof("associate user to business success with businessID: %d", lastInsertBusinessID)
	}

	b.logger.Infof("business added successfully with id: %d", lastInsertBusinessID)
	return lastInsertBusinessID, addressID, nil
}

func (b *businessEngine) addBusinessInfo(businessName string, phone string, website string) (int, error) {
	// insert business
	var lastInsertBusinessID int
	err := b.sql.QueryRow(insertBusinessSQL, businessName, phone, website).Scan(&lastInsertBusinessID)
	if err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}
	return lastInsertBusinessID, nil
}

func (b *businessEngine) addBusinessAddress(
	street string,
	city string,
	postalCode string,
	state string,
	businessID int,
) (int, error) {

	geoAddress := fmt.Sprintf("%s,%s,%s", street, city, state)
	resp, err := clients.GetLatLong(geoAddress)
	if err != nil {
		return 0, err
	}

	// insert to address table
	var addressID int
	err = b.sql.QueryRow(insertBusinessAddressSQL,
		businessID,
		street,
		city,
		postalCode,
		state,
		shared.CountryID,
		resp.Lat,
		resp.Lon,
	).Scan(&addressID)
	if err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}

	b.logger.Infof("successfully added address with ID: %d for business: %d", addressID, businessID)
	return addressID, nil
}

func (b *businessEngine) addBusinessHours(days []shared.Hours, businessID int) error {
	business, err := b.GetBusinessFromID(businessID)
	if err != nil {
		return err
	}

	if business.Name == "" {
		return helper.BusinessError{Message: fmt.Sprintf("business with id %d does not exist", businessID)}
	}

	for _, day := range days {
		if day.Day != "" {
			if day.OpenTimeSessionOne != "" && day.CloseTimeSessionOne != "" {
				if err := b.addHours(day.Day, day.OpenTimeSessionOne, day.CloseTimeSessionOne, businessID); err != nil {
					return err
				}
			}

			if day.OpenTimeSessionTwo != "" && day.CloseTimeSessionTwo != "" {
				if err := b.addHours(day.Day, day.OpenTimeSessionTwo, day.CloseTimeSessionTwo, businessID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (b *businessEngine) addHours(day string, openTime string, closeTime string, businessID int) error {
	addBusinessHoursSQL := "INSERT INTO business_hours(business_id,day,open_time,close_time) " +
		"VALUES($1,$2,$3,$4);"

	rows, err := b.sql.Query(addBusinessHoursSQL, businessID, day, openTime, closeTime)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	b.logger.Infof("add hours succesfull for businessID:%d", businessID)
	return nil
}

func (b *businessEngine) addBusinessCuisine(cuisines []string, businessID int) error {
	business, err := b.GetBusinessFromID(businessID)
	if err != nil {
		return err
	}

	if business.Name == "" {
		return helper.BusinessError{Message: fmt.Sprintf("business with id %d does not exist", businessID)}
	}

	for _, cuisine := range cuisines {
		if err := b.addCuisine(businessID, cuisine); err != nil {
			return err
		}
	}

	return nil
}

func (b *businessEngine) addCuisine(businessID int, cuisine string) error {
	addBusinessCuisineSQL := "INSERT INTO business_cuisine(business_id,cuisine) " +
		"VALUES($1,$2);"

	rows, err := b.sql.Query(addBusinessCuisineSQL, businessID, cuisine)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	b.logger.Infof("add cuisine successful for businessID:%d", businessID)
	return nil
}

func (b *businessEngine) associateUserToBusiness(businessID int, userID int) error {
	associateUserToBusinessSQL := "INSERT INTO user_to_business(business_id,user_id) " +
		"VALUES($1,$2);"

	rows, err := b.sql.Query(associateUserToBusinessSQL, businessID, userID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	b.logger.Infof("associateUserToBusiness successful for businessID:%d", businessID)
	return nil
}
