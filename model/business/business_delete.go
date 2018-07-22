package business

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

func (b *businessEngine) DeleteBusinessCuisineFromID(businessID int) error {
	sqlStatement := `DELETE FROM business_cuisine WHERE business_id = $1;`
	b.logger.Infof("deleting business_cuisine with query: %s and business_id: %d", sqlStatement, businessID)

	_, err := b.sql.Exec(sqlStatement, businessID)
	return err
}

func (b *businessEngine) DeleteBusinessHoursFromID(businessID int) error {
	sqlStatement := `DELETE FROM business_hours WHERE business_id = $1;`
	b.logger.Infof("deleting business_hours with query: %s and business_id: %d", sqlStatement, businessID)

	_, err := b.sql.Exec(sqlStatement, businessID)
	return err
}
