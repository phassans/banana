package business

func (b *businessEngine) deleteBusinessAddressFromID(addressID int) error {
	sqlStatement := `DELETE FROM business_address WHERE address_id = $1;`
	b.logger.Info().Msgf("deleting address with query: %s and business_id: %d", sqlStatement, addressID)

	_, err := b.sql.Exec(sqlStatement, addressID)
	return err
}

func (b *businessEngine) deleteBusinessFromID(businessID int) error {
	sqlStatement := `DELETE FROM business WHERE business_id = $1;`
	b.logger.Info().Msgf("deleting business with query: %s and business_id: %d", sqlStatement, businessID)

	_, err := b.sql.Exec(sqlStatement, businessID)
	return err
}

func (b *businessEngine) deleteBusinessCuisineFromID(businessID int) error {
	sqlStatement := `DELETE FROM business_cuisine WHERE business_id = $1;`
	b.logger.Info().Msgf("deleting business_cuisine with query: %s and business_id: %d", sqlStatement, businessID)

	_, err := b.sql.Exec(sqlStatement, businessID)
	return err
}

func (b *businessEngine) deleteBusinessHoursFromID(businessID int) error {
	sqlStatement := `DELETE FROM business_hours WHERE business_id = $1;`
	b.logger.Info().Msgf("deleting business_hours with query: %s and business_id: %d", sqlStatement, businessID)

	_, err := b.sql.Exec(sqlStatement, businessID)
	return err
}
