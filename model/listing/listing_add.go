package listing

import (
	"fmt"
	"strings"
	"time"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
	"github.com/rs/xlog"
)

func (l *listingEngine) AddListing(listing *shared.Listing) (int, error) {
	business, err := l.businessEngine.GetBusinessFromID(listing.BusinessID)
	if err != nil {
		return 0, err
	}

	if business.Name == "" {
		return 0, helper.BusinessError{Message: fmt.Sprintf("business with id %d does not exist", listing.BusinessID)}
	}

	var listingID int
	const insertListingSQL = "INSERT INTO listing(business_id, title, old_price, new_price, discount, description," +
		"start_date, start_time, end_time, multiple_days, end_date, recurring, recurring_end_date, listing_type, listing_create_date) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) returning listing_id"

	err = l.sql.QueryRow(insertListingSQL,
		listing.BusinessID,
		listing.Title,
		listing.OldPrice,
		listing.NewPrice,
		listing.Discount,
		listing.Description,
		listing.StartDate,
		listing.StartTime,
		listing.EndTime,
		listing.MultipleDays,
		shared.NewNullString(listing.EndDate),
		listing.Recurring,
		shared.NewNullString(listing.RecurringEndDate),
		listing.Type,
		time.Now()).
		Scan(&listingID)

	if err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}
	listing.ListingID = listingID

	// add listing image
	if listing.ImageLink != "" {
		if err := l.AddListingImage(listingID, listing.ImageLink); err != nil {
			return 0, err
		}
	} else {
		xlog.Info("no image link")
	}

	if listing.Recurring {
		for _, day := range listing.RecurringDays {
			if err := l.AddRecurring(listingID, day); err != nil {
				return 0, err
			}
		}
	}

	if len(listing.DietaryRestrictions) > 0 {
		for _, restriction := range listing.DietaryRestrictions {
			if err := l.AddDietaryRestriction(listingID, restriction); err != nil {
				return 0, err
			}
		}
	}

	// insert into listing_date
	if err := l.AddListingDates(listing); err != nil {
		return 0, err
	}

	l.logger.Infof("successfully added a listing %s for business: %s with listingId: %d", listing.Title, business.Name, listingID)

	return listingID, nil
}

func (l *listingEngine) AddListingDates(listing *shared.Listing) error {
	// current listing date
	listings := []shared.ListingDate{
		{ListingID: listing.ListingID, ListingDate: listing.StartDate, StartTime: listing.StartTime, EndTime: listing.EndTime},
	}

	listingDate, err := time.Parse(shared.DateFormat, strings.Split(listing.StartDate, "T")[0])
	if err != nil {
		return err
	}

	if listing.MultipleDays {
		listingEndDate, err := time.Parse(shared.DateFormat, strings.Split(listing.EndDate, "T")[0])
		if err != nil {
			return err
		}
		// difference b/w days
		days := listingEndDate.Sub(listingDate).Hours() / 24
		curDate := listingDate
		for i := 1; i <= int(days); i++ {
			var lDate shared.ListingDate
			nextDate := curDate.Add(time.Hour * 24)
			year, month, day := nextDate.Date()

			next := fmt.Sprintf("%d/%d/%d", int(month), day, year)
			lDate = shared.ListingDate{ListingID: listing.ListingID, ListingDate: next, StartTime: listing.StartTime, EndTime: listing.EndTime}
			listings = append(listings, lDate)

			curDate = nextDate
		}
	}

	if listing.Recurring {
		listingRecurringDate, err := time.Parse(shared.DateFormat, strings.Split(listing.RecurringEndDate, "T")[0])
		if err != nil {
			return err
		}
		// difference b/w days
		days := listingRecurringDate.Sub(listingDate).Hours() / 24
		curDate := listingDate
		for i := 1; i < int(days); i++ {
			var lDate shared.ListingDate
			nextDate := curDate.Add(time.Hour * 24)
			year, month, day := nextDate.Date()
			for _, recurringDay := range listing.RecurringDays {
				if shared.DayMap[recurringDay] == int(nextDate.Weekday()) {
					next := fmt.Sprintf("%d/%d/%d", int(month), day, year)
					lDate = shared.ListingDate{ListingID: listing.ListingID, ListingDate: next, StartTime: listing.StartTime, EndTime: listing.EndTime}
					listings = append(listings, lDate)
				}
			}
			curDate = nextDate
		}
	}

	for _, listing := range listings {
		err := l.InsertListingDate(listing)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *listingEngine) InsertListingDate(lDate shared.ListingDate) error {
	addListingDietRestrictionSQL := "INSERT INTO listing_date(listing_id,listing_date,start_time,end_time) " +
		"VALUES($1,$2,$3,$4);"

	rows, err := l.sql.Query(addListingDietRestrictionSQL, lDate.ListingID, lDate.ListingDate, lDate.StartTime, lDate.EndTime)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	l.logger.Infof("InsertListingDate successful for listing:%d", lDate.ListingID)
	return nil
}

func (l *listingEngine) AddListingImage(listingID int, imageLink string) error {
	addListingRecurringSQL := "INSERT INTO listing_image(listing_id,path) " +
		"VALUES($1,$2);"

	rows, err := l.sql.Query(addListingRecurringSQL, listingID, imageLink)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	l.logger.Infof("add AddListingImage successful for listing:%d", listingID)
	return nil
}

func (l *listingEngine) AddRecurring(listingID int, day string) error {
	addListingRecurringSQL := "INSERT INTO listing_recurring(listing_id,day) " +
		"VALUES($1,$2);"

	rows, err := l.sql.Query(addListingRecurringSQL, listingID, day)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	l.logger.Infof("add recurring successful for listing:%d", listingID)
	return nil
}

func (l *listingEngine) AddDietaryRestriction(listingID int, restriction string) error {
	addListingDietRestrictionSQL := "INSERT INTO listing_dietary_restrictions(listing_id,restriction) " +
		"VALUES($1,$2);"

	rows, err := l.sql.Query(addListingDietRestrictionSQL, listingID, restriction)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	l.logger.Infof("add listing_dietary_restrictions successful for listing:%d", listingID)
	return nil
}
