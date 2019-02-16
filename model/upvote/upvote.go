package upvote

import (
	"database/sql"
	"time"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/model/listing"
	"github.com/rs/zerolog"
)

type (
	upvoteEngine struct {
		sql           *sql.DB
		logger        zerolog.Logger
		listingEngine listing.ListingEngine
	}

	// UpvoteEngine interface which holds all methods
	UpvoteEngine interface {
		AddUpVote(phoneID string, listingID int) (int, error)
		DeleteUpVote(phoneID string, listingID int) (int, error)
		GetUpVotes(listingID int) (int, error)
	}
)

// NewUpvoteEngine returns an instance of upvoteEngine
func NewUpvoteEngine(psql *sql.DB, logger zerolog.Logger, listingEngine listing.ListingEngine) UpvoteEngine {
	return &upvoteEngine{psql, logger, listingEngine}
}

func (u *upvoteEngine) AddUpVote(phoneID string, listingID int) (int, error) {
	upvoteIDOld, err := u.GetUpVoteByPhoneID(phoneID, listingID)
	if err != nil {
		return 0, err
	}
	if upvoteIDOld != 0 {
		count, err := u.GetUpVotes(listingID)
		if err != nil {
			return 0, err
		}

		u.logger.Info().Msgf("already upvoted with ID: %d", upvoteIDOld)
		return count, nil
	}

	listing, err := u.listingEngine.GetListingByID(listingID, 0, 0)
	if err != nil {
		return 0, err
	}

	if listing.ListingID == 0 {
		return 0, helper.ListingDoesNotExist{ListingID: listingID}
	}

	var upvoteID int
	err = u.sql.QueryRow("INSERT INTO upvotes(phone_id,listing_id,upvote_date) "+
		"VALUES($1,$2,$3) returning upvote_id;",
		phoneID, listingID, time.Now()).Scan(&upvoteID)
	if err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}

	count, err := u.GetUpVotes(listingID)
	if err != nil {
		return 0, err
	}
	u.logger.Info().Msgf("successfully upvoted with ID: %d", upvoteID)
	return count, nil
}

func (u *upvoteEngine) DeleteUpVote(phoneID string, listingID int) (int, error) {
	sqlStatement := `DELETE FROM upvotes WHERE phone_id = $1 AND listing_id = $2;`
	u.logger.Info().Msgf("down voting with query: %s and listing: %d", sqlStatement, listingID)

	_, err := u.sql.Exec(sqlStatement, phoneID, listingID)

	count, err := u.GetUpVotes(listingID)
	if err != nil {
		return 0, err
	}

	return count, err
}

func (u *upvoteEngine) GetUpVotes(listingID int) (int, error) {
	rows, err := u.sql.Query("SELECT upvote_id FROM upvotes WHERE listing_id = $1", listingID)
	if err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}

	return count, nil
}

func (u *upvoteEngine) GetUpVoteByPhoneID(phoneID string, listingID int) (int, error) {
	var upvoteID int
	rows := u.sql.QueryRow("SELECT upvote_id FROM upvotes WHERE phone_id = $1 AND listing_id = $2", phoneID, listingID)
	err := rows.Scan(&upvoteID)

	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, helper.DatabaseError{DBError: err.Error()}
	}

	return upvoteID, nil
}
