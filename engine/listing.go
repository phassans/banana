package engine

import (
	"database/sql"

	"github.com/banana/helper"
	"github.com/rs/xlog"
)

type listingEngine struct {
	sql    *sql.DB
	logger xlog.Logger
}

type ListingEngine interface {
	CreateBusiness(businessName string) error
	GetBusinessIDFromName(businessName string) ([]int, error)
}

func NewListingEngine(psql *sql.DB, logger xlog.Logger) ListingEngine {
	return &listingEngine{psql, logger}
}

func (l *listingEngine) CreateBusiness(businessName string) error {
	// check for unique business name
	businessID, err := l.GetBusinessIDFromName(businessName)
	if err != nil {
		return err
	}

	switch len(businessID) {
	case 0:
		var lastInsertId int
		err := l.sql.QueryRow("INSERT INTO business(name) VALUES($1) returning business_id;", businessName).Scan(&lastInsertId)
		if err != nil {
			return helper.DatabaseError{DBError: err.Error()}
		}
		l.logger.Infof("last inserted id: %d", lastInsertId)
	default:
		return helper.DuplicateEntity{BusinessName: businessName}
	}

	return nil
}

func (l *listingEngine) GetBusinessIDFromName(businessName string) ([]int, error) {
	rows, err := l.sql.Query("SELECT business_id FROM business where name = $1;", businessName)
	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, helper.DatabaseError{DBError: err.Error()}
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	return ids, nil
}
