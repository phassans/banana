package business

import (
	"database/sql"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/model/user"
	"github.com/phassans/banana/shared"
	"github.com/rs/xlog"
)

type businessEngine struct {
	sql        *sql.DB
	logger     xlog.Logger
	userEngine user.UserEngine
}

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
		hoursInfo []shared.Hours,
		cuisine []string,
		UserId int,
	) (int, int, error)
	AddBusinessAddress(
		street string,
		city string,
		postalCode string,
		state string,
		businessID int,
	) (int, error)
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

	//Edit
	BusinessEdit(
		businessName string,
		phone string,
		website string,
		street string,
		city string,
		postalCode string,
		state string,
		hoursInfo []shared.Hours,
		cuisine []string,
		businessID int,
		addressID int,
	) (int, error)
}

func NewBusinessEngine(psql *sql.DB, logger xlog.Logger, userEngine user.UserEngine) BusinessEngine {
	return &businessEngine{psql, logger, userEngine}
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
	rows, err := l.sql.Query("SELECT street, city, postal_code, state, business_id, address_id FROM address where business_id = $1;", businessID)
	if err != nil {
		return shared.BusinessAddress{}, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var businessAddress shared.BusinessAddress
	if rows.Next() {
		err := rows.Scan(
			&businessAddress.Street,
			&businessAddress.City,
			&businessAddress.PostalCode,
			&businessAddress.State,
			&businessAddress.BusinessID,
			&businessAddress.AddressID,
		)
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

func (l *businessEngine) GetBusinessCuisineFromID(businessID int) (shared.BusinessCuisine, error) {
	rows, err := l.sql.Query("SELECT cuisine FROM business_cuisine where business_id = $1;", businessID)
	if err != nil {
		return shared.BusinessCuisine{}, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var businessCuisines []string
	for rows.Next() {
		var businessCuisine string
		err := rows.Scan(&businessCuisine)
		if err != nil {
			return shared.BusinessCuisine{}, helper.DatabaseError{DBError: err.Error()}
		}
		businessCuisines = append(businessCuisines, businessCuisine)
	}

	if err = rows.Err(); err != nil {
		return shared.BusinessCuisine{}, helper.DatabaseError{DBError: err.Error()}
	}

	return shared.BusinessCuisine{Cuisine: businessCuisines}, nil
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

	businessCuisine, err := l.GetBusinessCuisineFromID(businessID)
	if err != nil {
		return shared.BusinessInfo{}, err
	}

	return shared.BusinessInfo{Business: business, BusinessAddress: businessAddress, Hours: businessHours, BusinessCuisine: businessCuisine}, nil

}
