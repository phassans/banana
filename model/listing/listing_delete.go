package listing

import "github.com/phassans/banana/helper"

func (f *listingEngine) DeleteListing(listingID int) error {

	listingInfo, err := f.GetListingByID(listingID, 0, 0)
	if err != nil {
		return nil
	}

	if listingInfo.ListingID == 0 {
		return helper.ListingDoesNotExist{ListingID: listingID}
	}

	if err := f.deleteListingImage(listingID); err != nil {
		return nil
	}

	if err := f.deleteListingDietaryRestriction(listingID); err != nil {
		return nil
	}

	if err := f.deleteListingDate(listingID); err != nil {
		return nil
	}

	if err := f.deleteListingRecurring(listingID); err != nil {
		return nil
	}

	if err := f.deleteListing(listingID); err != nil {
		return nil
	}

	f.logger.Infof("successfully delete listing: %d", listingID)
	return nil
}

func (f *listingEngine) deleteListing(listingID int) error {
	sqlStatement := `DELETE FROM listing WHERE listing_id = $1;`
	f.logger.Infof("deleting listing with query: %s and listing: %d", sqlStatement, listingID)

	_, err := f.sql.Exec(sqlStatement, listingID)
	return err
}

func (f *listingEngine) deleteListingDate(listingID int) error {
	sqlStatement := `DELETE FROM listing_date WHERE listing_id = $1;`
	f.logger.Infof("deleting listing_date with query: %s and listing: %d", sqlStatement, listingID)

	_, err := f.sql.Exec(sqlStatement, listingID)
	return err
}

func (f *listingEngine) deleteListingRecurring(listingID int) error {
	sqlStatement := `DELETE FROM listing_recurring WHERE listing_id = $1;`
	f.logger.Infof("deleting listing_recurring with query: %s and listing: %d", sqlStatement, listingID)

	_, err := f.sql.Exec(sqlStatement, listingID)
	return err
}

func (f *listingEngine) deleteListingDietaryRestriction(listingID int) error {
	sqlStatement := `DELETE FROM listing_dietary_restrictions WHERE listing_id = $1;`
	f.logger.Infof("deleting listing_dietary_restrictions with query: %s and listing: %d", sqlStatement, listingID)

	_, err := f.sql.Exec(sqlStatement, listingID)
	return err
}

func (f *listingEngine) deleteListingImage(listingID int) error {
	sqlStatement := `DELETE FROM listing_image WHERE listing_id = $1;`
	f.logger.Infof("deleting listing_image with query: %s and listing: %d", sqlStatement, listingID)

	_, err := f.sql.Exec(sqlStatement, listingID)
	return err
}
