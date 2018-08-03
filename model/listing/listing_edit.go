package listing

import (
	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
	"github.com/rs/xlog"
)

func (l *listingEngine) ListingEdit(listing *shared.Listing) error {
	listingInfo, err := l.GetListingByID(listing.ListingID, listing.BusinessID, 0)
	if err != nil {
		return err
	}

	if listingInfo.ListingID == 0 {
		return helper.ListingDoesNotExist{ListingID: listing.ListingID}
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

	if len(listing.DietaryRestrictions) > 0 {
		for _, restriction := range listing.DietaryRestrictions {
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
	SET title = $1, old_price = $2, new_price = $3, discount = $4, discount_description = $5, description = $6, start_date = $7, start_time = $8, 
	end_time = $9, multiple_days = $10, end_date = $11, recurring = $12, recurring_end_date = $13
	WHERE listing_id = $14 AND business_id = $15 AND listing_type = $16;`

	_, err := l.sql.Exec(
		updateListingInfoSQL,
		listing.Title,
		listing.OldPrice,
		listing.NewPrice,
		listing.Discount,
		listing.DiscountDescription,
		listing.Description,
		listing.StartDate,
		listing.StartTime,
		listing.EndTime,
		listing.MultipleDays,
		shared.NewNullString(listing.EndDate),
		listing.Recurring,
		shared.NewNullString(listing.RecurringEndDate),
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
