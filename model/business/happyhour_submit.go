package business

import (
	"time"

	"github.com/phassans/banana/helper"
)

func (b *businessEngine) SubmitHappyHour(PhoneID string, Name string, Email string, BusinessOwner bool, Restaurant string, City string, Description string) (int, error) {
	// insert happyhour
	var lastInsertHHID int
	err := b.sql.QueryRow("INSERT INTO happyhour(phone_id,name,email,business_owner,restaurant,city,description,submission_date) "+
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8) returning hh_id;", PhoneID, Name, Email, BusinessOwner, Restaurant, City, Description, time.Now()).Scan(&lastInsertHHID)
	if err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}
	return lastInsertHHID, nil
}

func (b *businessEngine) SubmitHappyHourImages(happyHourID int, imageName string) (int, error) {
	// insert happyhour images
	var lastInsertHHID int
	err := b.sql.QueryRow("INSERT INTO happyhour_images(hh_id,image_name,submission_date) "+
		"VALUES($1,$2,$3) returning hh_id;", happyHourID, imageName, time.Now()).Scan(&lastInsertHHID)
	if err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}
	return lastInsertHHID, nil
}
