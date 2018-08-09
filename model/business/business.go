package business

import (
	"database/sql"

	"bytes"

	"strings"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/model/user"
	"github.com/phassans/banana/shared"
	"github.com/rs/zerolog"
)

type businessEngine struct {
	sql        *sql.DB
	logger     zerolog.Logger
	userEngine user.UserEngine
}

// BusinessEngine an interface for business operations
type BusinessEngine interface {
	// AddBusiness function to add a business
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
		UserID int,
	) (int, int, error)

	// GetBusinessFromID returns business from ID
	GetBusinessFromID(businessID int) (shared.Business, error)

	// GetBusinessInfo returns business info from ID
	GetBusinessInfo(businessID int) (shared.BusinessInfo, error)

	// GetAllBusiness returns all business
	GetAllBusiness() ([]shared.Business, error)

	// BusinessEdit is to edit business info
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

// NewBusinessEngine returns an instance of businessEngine
func NewBusinessEngine(psql *sql.DB, logger zerolog.Logger, userEngine user.UserEngine) BusinessEngine {
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
		err = rows.Scan(&business.BusinessID, &business.Name, &business.Phone, &business.Website)
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

func (b *businessEngine) GetBusinessIDFromName(businessName string) (int, error) {
	rows, err := b.sql.Query("SELECT business_id FROM business where name = $1;", businessName)
	if err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var id int
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return 0, helper.DatabaseError{DBError: err.Error()}
		}
	}

	if err = rows.Err(); err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}

	return id, nil
}

func (b *businessEngine) GetBusinessFromID(businessID int) (shared.Business, error) {
	rows, err := b.sql.Query("SELECT business_id, name, phone, website FROM business where business_id = $1;", businessID)
	if err != nil {
		return shared.Business{}, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var business shared.Business
	if rows.Next() {
		err = rows.Scan(&business.BusinessID, &business.Name, &business.Phone, &business.Website)
		if err != nil {
			return shared.Business{}, helper.DatabaseError{DBError: err.Error()}
		}
	}

	if err = rows.Err(); err != nil {
		return shared.Business{}, helper.DatabaseError{DBError: err.Error()}
	}

	return business, nil
}

func (b *businessEngine) GetBusinessInfo(businessID int) (shared.BusinessInfo, error) {
	business, err := b.GetBusinessFromID(businessID)
	if err != nil {
		return shared.BusinessInfo{}, err
	}

	businessAddress, err := b.getBusinessAddressFromID(businessID)
	if err != nil {
		return shared.BusinessInfo{}, err
	}

	businessHours, err := b.getBusinessHoursFromID(businessID)
	if err != nil {
		return shared.BusinessInfo{}, err
	}

	businessHoursFormatted, err := getBusinessHoursFormatted(businessHours)
	if err != nil {
		return shared.BusinessInfo{}, err
	}

	businessCuisine, err := b.getBusinessCuisineFromID(businessID)
	if err != nil {
		return shared.BusinessInfo{}, err
	}

	return shared.BusinessInfo{Business: business, BusinessAddress: businessAddress, BusinessCuisine: businessCuisine, HoursFormatted: businessHoursFormatted}, nil

}

var days = []string{
	"monday",
	"tuesday",
	"wednesday",
	"thursday",
	"friday",
	"saturday",
	"sunday",
}

func getBusinessHoursFormatted(bHours []shared.Bhour) ([]string, error) {
	var bHoursFormatted []string
	bMap := make(map[string]string)

	for _, hour := range bHours {
		var buffer bytes.Buffer
		// determine startTime in format
		oTime, err := shared.GetTimeIn12HourFormat(hour.OpenTime)
		if err != nil {
			return nil, err
		}

		// determine startTime in format
		cTime, err := shared.GetTimeIn12HourFormat(hour.CloseTime)
		if err != nil {
			return nil, err
		}

		if val, ok := bMap[hour.Day]; ok {
			buffer.WriteString(", ")
			buffer.WriteString(oTime)
			buffer.WriteString("-")
			buffer.WriteString(cTime)
			bMap[hour.Day] = val + buffer.String()
		} else {
			buffer.WriteString(strings.Title(hour.Day[0:3]))
			buffer.WriteString(": ")
			buffer.WriteString(oTime)
			buffer.WriteString("-")
			buffer.WriteString(cTime)
			bMap[hour.Day] = buffer.String()
		}
	}

	for _, val := range days {
		bHoursFormatted = append(bHoursFormatted, bMap[val])
	}

	return bHoursFormatted, nil
}

func (b *businessEngine) getBusinessAddressFromID(businessID int) (shared.BusinessAddress, error) {
	rows, err := b.sql.Query("SELECT street, city, postal_code, state, business_id, address_id FROM business_address where business_id = $1;", businessID)
	if err != nil {
		return shared.BusinessAddress{}, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var businessAddress shared.BusinessAddress
	if rows.Next() {
		err = rows.Scan(
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

func (b *businessEngine) getBusinessHoursFromID(businessID int) ([]shared.Bhour, error) {
	rows, err := b.sql.Query("SELECT day, open_time, close_time FROM business_hours where business_id = $1;", businessID)
	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var businessHours []shared.Bhour
	for rows.Next() {
		var businessHour shared.Bhour
		err = rows.Scan(&businessHour.Day, &businessHour.OpenTime, &businessHour.CloseTime)
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

func (b *businessEngine) getBusinessCuisineFromID(businessID int) (shared.BusinessCuisine, error) {
	rows, err := b.sql.Query("SELECT cuisine FROM business_cuisine where business_id = $1;", businessID)
	if err != nil {
		return shared.BusinessCuisine{}, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var businessCuisines []string
	for rows.Next() {
		var businessCuisine string
		err = rows.Scan(&businessCuisine)
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
