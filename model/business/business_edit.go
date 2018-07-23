package business

import (
	"fmt"

	"github.com/phassans/banana/clients"
	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
	"github.com/rs/xlog"
)

func (b *businessEngine) BusinessEdit(
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
) (int, error) {
	// check business name unique
	businessInfo, err := b.GetBusinessFromID(businessID)
	if err != nil {
		return 0, err
	}

	if businessInfo.BusinessID == 0 {
		return 0, helper.BusinessDoesNotExist{BusinessName: businessName}
	}

	// edit business
	if err := b.editBusinessInfo(businessName, phone, website, businessID); err != nil {
		return 0, err
	}
	xlog.Infof("editBusinessInfo success with businessID: %d", businessID)

	// edit business address
	if err := b.editBusinessAddress(street, city, postalCode, state, businessID, addressID); err != nil {
		return 0, err
	}
	xlog.Infof("editBusinessAddress success with businessID: %d", businessID)

	// edit business cuisine
	if err := b.editBusinessCuisine(cuisine, businessID); err != nil {
		return 0, err
	}
	xlog.Infof("editBusinessCuisine success with businessID: %d", businessID)

	// edit business hours
	if err := b.editBusinessHours(hoursInfo, businessID); err != nil {
		return 0, err
	}
	xlog.Infof("editBusinessHours success with businessID: %d", businessID)

	return 0, nil
}

func (b *businessEngine) editBusinessInfo(businessName string, phone string, website string, businessID int) error {

	updateBusinessSQL := `
	UPDATE business
	SET name = $1, phone = $2, website = $3
	WHERE business_id = $4;`

	_, err := b.sql.Exec(updateBusinessSQL, businessName, phone, website, businessID)
	if err != nil {
		return err
	}

	return nil
}

func (b *businessEngine) editBusinessAddress(street string, city string, postalCode string, state string, businessID int, addressID int) error {

	updateBusinessAddressSQL := `
	UPDATE address
	SET street = $1, city = $2, postal_code = $3, state =$4
	WHERE business_id = $5 AND address_id = $6;`

	_, err := b.sql.Exec(updateBusinessAddressSQL, street, city, postalCode, state, businessID, addressID)
	if err != nil {
		return err
	}

	// edit business address
	geoAddress := fmt.Sprintf("%s,%s,%s", street, city, state)
	if err := b.editBusinessAddressGeo(geoAddress, businessID, addressID); err != nil {
		return err
	}

	return nil
}

func (b *businessEngine) editBusinessAddressGeo(geoAddress string, businessID int, addressID int) error {

	// edit lat, long to database
	resp, err := clients.GetLatLong(geoAddress)
	if err != nil {
		return err
	}

	updateBusinessAddressGeoSQL := `
	UPDATE address_geo
	SET latitude = $1, longitude = $2
	WHERE address_id = $3 AND business_id=$4;`

	_, err = b.sql.Exec(updateBusinessAddressGeoSQL, resp.Lat, resp.Lon, addressID, businessID)
	if err != nil {
		return err
	}

	return nil
}

func (b *businessEngine) editBusinessHours(bhours []shared.Hours, businessID int) error {

	// delete and add cuisine
	if err := b.DeleteBusinessHoursFromID(businessID); err != nil {
		return err
	}

	// add cuisine
	if err := b.AddBusinessHours(bhours, businessID); err != nil {
		return err
	}

	return nil
}

func (b *businessEngine) editBusinessCuisine(cuisines []string, businessID int) error {

	// delete and add cuisine
	if err := b.DeleteBusinessCuisineFromID(businessID); err != nil {
		return err
	}

	// add cuisine
	if err := b.AddBusinessCuisine(cuisines, businessID); err != nil {
		return err
	}

	return nil
}
