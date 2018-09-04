package prefernce

import (
	"database/sql"
	"fmt"

	"github.com/phassans/banana/shared"

	"github.com/phassans/banana/helper"
	"github.com/rs/zerolog"
)

type preferenceEngine struct {
	sql    *sql.DB
	logger zerolog.Logger
}

// PreferenceEngine an interface for preference operations
type PreferenceEngine interface {
	PreferenceAdd(phoneID string, cuisine string) error
	PreferenceDelete(phoneID string, cuisine string) error
	PreferenceAll(phoneID string) ([]shared.Preference, error)
}

// NewPreferenceEngine returns an instance of notificationEngine
func NewPreferenceEngine(psql *sql.DB, logger zerolog.Logger) PreferenceEngine {
	return &preferenceEngine{psql, logger}
}

func (p *preferenceEngine) PreferenceAdd(phoneID string, cuisine string) error {
	preference, err := p.CheckIfPreferenceExists(phoneID, cuisine)
	if err != nil {
		return err
	}

	if preference.PreferenceID != 0 {
		return nil
	}

	addPreferenceSQL := "INSERT INTO preferences(phone_id,cuisine) " +
		"VALUES($1,$2);"

	rows, err := p.sql.Query(addPreferenceSQL, phoneID, cuisine)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	p.logger.Info().Msgf("add preferences successful for phoneID:%d", phoneID)

	return nil
}

func (p *preferenceEngine) PreferenceDelete(phoneID string, cuisine string) error {
	sqlStatement := `DELETE FROM preferences WHERE phone_id = $1 AND cuisine = $2;`
	p.logger.Info().Msgf("deleting preferences with query: %s and phoneID: %s and cuisine: %s", sqlStatement, phoneID, cuisine)

	_, err := p.sql.Exec(sqlStatement, phoneID, cuisine)
	return err
}

func (p *preferenceEngine) CheckIfPreferenceExists(phoneID string, cuisine string) (shared.Preference, error) {
	q := fmt.Sprintf("SELECT prefernce_id, phone_id, cuisine FROM preferences where phone_id = '%s' AND cuisine = '%s';", phoneID, cuisine)
	rows, err := p.sql.Query(q)
	if err != nil {
		return shared.Preference{}, helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	var preference shared.Preference
	if rows.Next() {
		err = rows.Scan(
			&preference.PreferenceID,
			&preference.PhoneID,
			&preference.Cuisine,
		)
		if err != nil {
			return shared.Preference{}, helper.DatabaseError{DBError: err.Error()}
		}
	}

	if err = rows.Err(); err != nil {
		return shared.Preference{}, helper.DatabaseError{DBError: err.Error()}
	}

	return preference, nil
}

func (p *preferenceEngine) PreferenceAll(phoneID string) ([]shared.Preference, error) {
	q := fmt.Sprintf("SELECT prefernce_id, phone_id, cuisine FROM preferences where phone_id = '%s';", phoneID)
	rows, err := p.sql.Query(q)
	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	var preferences []shared.Preference
	for rows.Next() {
		var preference shared.Preference
		err = rows.Scan(
			&preference.PreferenceID,
			&preference.PhoneID,
			&preference.Cuisine,
		)
		if err != nil {
			return nil, helper.DatabaseError{DBError: err.Error()}
		}

		preferences = append(preferences, preference)
	}

	if err = rows.Err(); err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	return preferences, nil
}
