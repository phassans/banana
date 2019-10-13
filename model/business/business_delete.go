package business

import (
	"bytes"
	"database/sql"
	"fmt"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/model/common"
	"github.com/phassans/banana/shared"
)

func (b *businessEngine) BusinessDelete(businessID int) error {

	listings, err := b.GetListingsByBusinessID(businessID)
	if err != nil {
		return err
	}

	if len(listings) > 0 {
		var lists string
		for _, l := range listings {
			lists += l.Title + ","
		}
		return helper.BusinessError{fmt.Sprintf("cannot delete business. Followings listings are tied - %s", lists[:len(lists)-1])}
	}

	if err := b.deleteBusinessHoursFromID(businessID); err != nil {
		return err
	}

	if err := b.deleteBusinessCuisineFromID(businessID); err != nil {
		return err
	}

	if err := b.deleteBusinessAddressFromBusinessID(businessID); err != nil {
		return err
	}

	if err := b.deleteBusinessFromID(businessID); err != nil {
		return err
	}

	return nil
}

func (b *businessEngine) deleteBusinessAddressFromBusinessID(businessID int) error {
	sqlStatement := `DELETE FROM business_address WHERE business_id = $1;`
	b.logger.Info().Msgf("deleting address with query: %s and business_id: %d", sqlStatement, businessID)

	_, err := b.sql.Exec(sqlStatement, businessID)
	return err
}

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

func (b *businessEngine) GetListingsByBusinessID(businessID int) ([]shared.Listing, error) {
	selectFields := fmt.Sprintf("%s, %s", common.ListingFields, common.ListingBusinessFields)

	var whereClause bytes.Buffer
	whereClause.WriteString(fmt.Sprintf(" WHERE listing.business_id = %d", businessID))
	query := fmt.Sprintf("%s %s %s;", selectFields, common.FromClauseBusinessListing, whereClause.String())

	//fmt.Println("GetListingsByBusinessID ", query)

	rows, err := b.sql.Query(query)
	if err != nil {
		return []shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	var listings []shared.Listing
	var sqlEndDate sql.NullString
	var sqlRecurringEndDate sql.NullString
	var sqlCreateDate sql.NullString
	for rows.Next() {
		var listing shared.Listing
		err = rows.Scan(
			&listing.Title,
			&listing.OldPrice,
			&listing.NewPrice,
			&listing.Discount,
			&listing.DiscountDescription,
			&listing.Description,
			&listing.StartDate,
			&sqlEndDate,
			&listing.StartTime,
			&listing.EndTime,
			&listing.MultipleDays,
			&listing.Recurring,
			&sqlRecurringEndDate,
			&listing.Type,
			&listing.BusinessID,
			&listing.ListingID,
			&sqlCreateDate,
			&listing.BusinessName,
		)
		if err != nil {
			return []shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
		}
		listing.EndDate = sqlEndDate.String
		listing.RecurringEndDate = sqlRecurringEndDate.String
		listing.ListingCreateDate = sqlCreateDate.String

		listings = append(listings, listing)
	}

	if err = rows.Err(); err != nil {
		return []shared.Listing{}, helper.DatabaseError{DBError: err.Error()}
	}

	return listings, nil
}
