package model

import (
	"database/sql"

	"fmt"
	"net/url"

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

	insertBusinessAddressSQL = "INSERT INTO address(street,city,postal_code,state,country_id,business_id) " +
		"VALUES($1,$2,$3,$4,$5,$6) returning address_id;"

	insertBusinessAddressGEOSQL = "INSERT INTO address_geo(address_id,business_id,latitude,longitude) " +
		"VALUES($1,$2,$3,$4) returning geo_id;"

	countryID = 1
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
	AddBusinessHours([]Days, int) error
	AddBusinessImage(businessName string, imagePath string)

	DeleteBusinessFromID(business_id int) error
	DeleteBusinessAddressFromID(business_id int) error
}

func NewBusinessEngine(psql *sql.DB, logger xlog.Logger) BusinessEngine {
	return &businessEngine{psql, logger}
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
) (int, error) {
	businessID, err := b.GetBusinessIDFromName(businessName)
	if err != nil {
		return 0, err
	}

	if businessID != 0 {
		return 0, helper.DuplicateEntity{BusinessName: businessName}
	}

	var lastInsertBusinessID int
	err = b.sql.QueryRow(insertBusinessSQL, businessName, phone, website).
		Scan(&lastInsertBusinessID)
	if err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}
	b.logger.Infof("business addded successfully with id: %d", lastInsertBusinessID)

	if err = b.AddBusinessAddress(street, city, postalCode, state, country, lastInsertBusinessID); err != nil {
		// cleanup
		go b.DeleteBusinessFromID(lastInsertBusinessID)
		return 0, nil
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
		state, countryID, businessID).
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

func (l *businessEngine) AddBusinessHours(days []Days, businessID int) error {
	businessName, err := l.GetBusinessFromID(businessID)
	if err != nil {
		return err
	}

	if businessName == "" {
		return helper.BusinessError{Message: fmt.Sprintf("business with id %d does not exist", businessID)}
	}

	for _, day := range days {
		switch day.(type) {
		case HoursMonday:
			request := day.(HoursMonday)
			if request.Monday {
				if err := l.AddHours("monday", request.MondayOpenTimeSessionOne, request.MondayCloseTimeSessionOne, businessID); err != nil {
					return err
				}

				if request.MondayOpenTimeSessionTwo != "" && request.MondayCloseTimeSessionTwo != "" {
					if err := l.AddHours("monday", request.MondayOpenTimeSessionTwo, request.MondayCloseTimeSessionTwo, businessID); err != nil {
						return err
					}
				}
			}
		case HoursTuesday:
			request := day.(HoursTuesday)
			if request.Tuesday {
				if err := l.AddHours("tuesday", request.TuesdayOpenTimeSessionOne, request.TuesdayCloseTimeSessionOne, businessID); err != nil {
					return err
				}

				if request.TuesdayOpenTimeSessionTwo != "" && request.TuesdayCloseTimeSessionTwo != "" {
					if err := l.AddHours("tuesday", request.TuesdayOpenTimeSessionTwo, request.TuesdayCloseTimeSessionTwo, businessID); err != nil {
						return err
					}
				}
			}
		case HoursWednesday:
			request := day.(HoursWednesday)
			if request.Wednesday {
				if err := l.AddHours("wednesday", request.WednesdayOpenTimeSessionOne, request.WednesdayCloseTimeSessionOne, businessID); err != nil {
					return err
				}

				if request.WednesdayOpenTimeSessionTwo != "" && request.WednesdayCloseTimeSessionTwo != "" {
					if err := l.AddHours("wednesday", request.WednesdayOpenTimeSessionTwo, request.WednesdayCloseTimeSessionTwo, businessID); err != nil {
						return err
					}
				}
			}
		case HoursThursday:
			request := day.(HoursThursday)
			if request.Thursday {
				if err := l.AddHours("thursday", request.ThursdayOpenTimeSessionOne, request.ThursdayCloseTimeSessionOne, businessID); err != nil {
					return err
				}

				if request.ThursdayOpenTimeSessionTwo != "" && request.ThursdayCloseTimeSessionTwo != "" {
					if err := l.AddHours("thursday", request.ThursdayOpenTimeSessionTwo, request.ThursdayCloseTimeSessionTwo, businessID); err != nil {
						return err
					}
				}
			}
		case HoursFriday:
			request := day.(HoursFriday)
			if request.Friday {
				if err := l.AddHours("friday", request.FridayOpenTimeSessionOne, request.FridayCloseTimeSessionOne, businessID); err != nil {
					return err
				}

				if request.FridayOpenTimeSessionTwo != "" && request.FridayCloseTimeSessionTwo != "" {
					if err := l.AddHours("friday", request.FridayOpenTimeSessionTwo, request.FridayCloseTimeSessionTwo, businessID); err != nil {
						return err
					}
				}
			}
		case HoursSaturday:
			request := day.(HoursSaturday)
			if request.Saturday {
				if err := l.AddHours("saturday", request.SaturdayOpenTimeSessionOne, request.SaturdayCloseTimeSessionOne, businessID); err != nil {
					return err
				}

				if request.SaturdayOpenTimeSessionTwo != "" && request.SaturdayCloseTimeSessionTwo != "" {
					if err := l.AddHours("saturday", request.SaturdayOpenTimeSessionTwo, request.SaturdayCloseTimeSessionTwo, businessID); err != nil {
						return err
					}
				}
			}
		case HoursSunday:
			request := day.(HoursSunday)
			if request.Sunday {
				if err := l.AddHours("sunday", request.SundayOpenTimeSessionOne, request.SundayCloseTimeSessionOne, businessID); err != nil {
					return err
				}

				if request.SundayOpenTimeSessionTwo != "" && request.SundayCloseTimeSessionTwo != "" {
					if err := l.AddHours("sunday", request.SundayOpenTimeSessionTwo, request.SundayCloseTimeSessionTwo, businessID); err != nil {
						return err
					}
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

func (l *businessEngine) GetBusinessFromID(businessID int) (string, error) {
	rows, err := l.sql.Query("SELECT name FROM business where business_id = $1;", businessID)
	if err != nil {
		return "", helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var businessName string
	if rows.Next() {
		err := rows.Scan(&businessName)
		if err != nil {
			return "", helper.DatabaseError{DBError: err.Error()}
		}
	}

	if err = rows.Err(); err != nil {
		return "", helper.DatabaseError{DBError: err.Error()}
	}

	return businessName, nil
}

func (b *businessEngine) DeleteBusinessAddressFromID(address_id int) error {
	sqlStatement := `DELETE FROM address WHERE address_id = $1;`
	b.logger.Infof("deleting address with query: %s and business_id: %d", sqlStatement, address_id)

	_, err := b.sql.Exec(sqlStatement, address_id)
	return err
}

func (b *businessEngine) DeleteBusinessFromID(business_id int) error {
	sqlStatement := `DELETE FROM business WHERE business_id = $1;`
	b.logger.Infof("deleting business with query: %s and business_id: %d", sqlStatement, business_id)

	_, err := b.sql.Exec(sqlStatement, business_id)
	return err
}
