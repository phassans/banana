package model

import (
	"database/sql"

	"github.com/pshassans/banana/helper"
	"github.com/rs/xlog"
)

type ownerEngine struct {
	sql            *sql.DB
	logger         xlog.Logger
	businessEngine BusinessEngine
}

type OwnerEngine interface {
	AddOwner(firstName string, lastName string, phone string, email string, businessName string) error
}

func NewOwnerEngine(psql *sql.DB, logger xlog.Logger, businessEngine BusinessEngine) OwnerEngine {
	return &ownerEngine{psql, logger, businessEngine}
}

func (l *ownerEngine) AddOwner(firstName string, lastName string, phone string, email string, businessName string) error {
	businessID, err := l.businessEngine.GetBusinessIDFromName(businessName)
	if err != nil {
		return err
	}

	if businessID == 0 {
		return helper.BusinessDoesNotExist{BusinessName: businessName}
	}

	var ownerID int

	err = l.sql.QueryRow("INSERT INTO owner(first_name,last_name,phone,email,business_id) "+
		"VALUES($1,$2,$3,$4,$5) returning owner_id;",
		firstName, lastName, phone, email, businessID).Scan(&ownerID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	l.logger.Infof("successfully added a user with ID: %d", ownerID)

	return nil
}
