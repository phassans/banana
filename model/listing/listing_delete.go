package listing

import "github.com/phassans/banana/helper"

func (f *listingEngine) DeleteListing(listingID int) error {

	listingInfo, err := f.GetListingByID(listingID, 0, 0)
	if err != nil {
		f.logger.Error().Msgf("GetListingByID returned with err: %s", err)
		return nil
	}

	if listingInfo.ListingID == 0 {
		return helper.ListingDoesNotExist{ListingID: listingID}
	}

	if err := f.deleteListingImage(listingID); err != nil {
		f.logger.Error().Msgf("deleteListingImage returned with err: %s", err)
		return nil
	}

	if err := f.deleteListingDietaryRestriction(listingID); err != nil {
		f.logger.Error().Msgf("deleteListingDietaryRestriction returned with err: %s", err)
		return nil
	}

	if err := f.deleteListingDate(listingID); err != nil {
		f.logger.Error().Msgf("deleteListingDate returned with err: %s", err)
		return nil
	}

	if err := f.deleteListingRecurring(listingID); err != nil {
		f.logger.Error().Msgf("deleteListingRecurring returned with err: %s", err)
		return nil
	}

	if err := f.deleteListing(listingID); err != nil {
		f.logger.Error().Msgf("deleteListing returned with err: %s", err)
		return nil
	}

	f.logger.Info().Msgf("successfully delete listing: %d", listingID)
	return nil
}

func (f *listingEngine) deleteListing(listingID int) error {
	sqlStatement := `DELETE FROM listing WHERE listing_id = $1;`
	f.logger.Info().Msgf("deleting listing with query: %s and listing: %d", sqlStatement, listingID)

	_, err := f.sql.Exec(sqlStatement, listingID)
	return err
}

func (f *listingEngine) deleteListingDate(listingID int) error {
	sqlStatement := `DELETE FROM listing_date WHERE listing_id = $1;`
	f.logger.Info().Msgf("deleting listing_date with query: %s and listing: %d", sqlStatement, listingID)

	_, err := f.sql.Exec(sqlStatement, listingID)
	return err
}

func (f *listingEngine) deleteListingRecurring(listingID int) error {
	sqlStatement := `DELETE FROM listing_recurring WHERE listing_id = $1;`
	f.logger.Info().Msgf("deleting listing_recurring with query: %s and listing: %d", sqlStatement, listingID)

	_, err := f.sql.Exec(sqlStatement, listingID)
	return err
}

func (f *listingEngine) deleteListingDietaryRestriction(listingID int) error {
	sqlStatement := `DELETE FROM listing_dietary_restrictions WHERE listing_id = $1;`
	f.logger.Info().Msgf("deleting listing_dietary_restrictions with query: %s and listing: %d", sqlStatement, listingID)

	_, err := f.sql.Exec(sqlStatement, listingID)
	return err
}

func (f *listingEngine) deleteListingImage(listingID int) error {
	sqlStatement := `DELETE FROM listing_image WHERE listing_id = $1;`
	f.logger.Info().Msgf("deleting listing_image with query: %s and listing: %d", sqlStatement, listingID)

	_, err := f.sql.Exec(sqlStatement, listingID)
	return err
}
