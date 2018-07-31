package listing

import (
	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
	"github.com/rs/xlog"
)

func (l *listingEngine) ListingEdit(listing *shared.Listing) error {
	listingInfo, err := l.GetListingByID(listing.ListingID, listing.BusinessID)
	if err != nil {
		return err
	}

	if listingInfo.ListingID == 0 {
		return helper.ListingDoesNotExist{ListingID: listing.ListingID}
	}

	if listing.RecurringEndDate == "" {
		listing.RecurringEndDate = "01/01/2000"
	}

	// edit listing info
	if err := l.editListingInfo(listing); err != nil {
		return err
	}
	xlog.Infof("editListingInfo success for listing: %d", listing.ListingID)

	// edit listing image
	if err := l.editListingImage(listing); err != nil {
		return err
	}
	xlog.Infof("editListingImage success for listing: %d", listing.ListingID)

	// edit recurring days
	if err := l.editRecurringDays(listing); err != nil {
		return err
	}
	xlog.Infof("editRecurringDays success for listing: %d", listing.ListingID)

	// edit dietary restriction
	if err := l.editDietaryRestriction(listing); err != nil {
		return err
	}
	xlog.Infof("editDietaryRestriction success for listing: %d", listing.ListingID)

	// edit listing dates
	if err := l.editListingDates(listing); err != nil {
		return err
	}
	xlog.Infof("editListingDates success for listing: %d", listing.ListingID)

	return nil
}

func (l *listingEngine) editRecurringDays(listing *shared.Listing) error {
	if err := l.deleteListingRecurring(listing.ListingID); err != nil {
		return err
	}

	if listing.Recurring {
		for _, day := range listing.RecurringDays {
			if err := l.AddRecurring(listing.ListingID, day); err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *listingEngine) editDietaryRestriction(listing *shared.Listing) error {
	if err := l.deleteListingDietaryRestriction(listing.ListingID); err != nil {
		return err
	}

	if len(listing.DietaryRestriction) > 0 {
		for _, restriction := range listing.DietaryRestriction {
			if err := l.AddDietaryRestriction(listing.ListingID, restriction); err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *listingEngine) editListingDates(listing *shared.Listing) error {
	if err := l.deleteListingDate(listing.ListingID); err != nil {
		return err
	}

	// insert into listing_date
	if err := l.AddListingDates(listing); err != nil {
		return err
	}

	return nil
}

func (l *listingEngine) editListingInfo(listing *shared.Listing) error {
	updateListingInfoSQL := `
	UPDATE listing
	SET title = $1, old_price = $2, new_price = $3, discount = $4, description = $5, start_date = $6, start_time = $7, 
	end_time = $8, multiple_days = $9, end_date = $10, recurring = $11, recurring_end_date = $12
	WHERE listing_id = $13 AND business_id = $14 AND listing_type = $15;`

	_, err := l.sql.Exec(
		updateListingInfoSQL,
		listing.Title,
		listing.OldPrice,
		listing.NewPrice,
		listing.Discount,
		listing.Description,
		listing.StartDate,
		listing.StartTime,
		listing.EndTime,
		listing.MultipleDays,
		listing.EndDate,
		listing.Recurring,
		listing.RecurringEndDate,
		listing.ListingID,
		listing.BusinessID,
		listing.Type,
	)
	if err != nil {
		return err
	}
	return nil
}

func (l *listingEngine) editListingImage(listing *shared.Listing) error {
	updateListingImageSQL := `
	UPDATE listing_image
	SET path = $1
	WHERE listing_id = $2`

	_, err := l.sql.Exec(
		updateListingImageSQL,
		listing.ImageLink,
		listing.ListingID,
	)
	if err != nil {
		return err
	}
	return nil
}
