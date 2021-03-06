package user

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"time"

	"fmt"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
	"github.com/rs/zerolog"
)

type userEngine struct {
	sql    *sql.DB
	logger zerolog.Logger
}

// UserEngine interface
type UserEngine interface {
	UserAdd(name string, email string, password string, phone string) error
	UserEdit(userID int, name string, email string, password string, phone string) error
	UserGet(userID int) (shared.BusinessUser, error)
	UserVerify(email string, password string) (shared.BusinessUser, error)

	ContactUs(phoneID string, uname string, email string, subject string, comments string) error
	ReportInaccurate(phoneID string, listing int) error
}

// NewUserEngine returns an instance of userEngine
func NewUserEngine(psql *sql.DB, logger zerolog.Logger) UserEngine {
	return &userEngine{psql, logger}
}

// signup flow
func (u *userEngine) UserAdd(name string, email string, password string, phone string) error {
	// check business name unique
	userID, err := u.CheckEmail(email)
	if err != nil {
		return err
	}

	if userID != 0 {
		return helper.UserError{Message: fmt.Sprintf("user with %s is already registered!", email)}
	}

	err = u.sql.QueryRow("INSERT INTO business_user(name,email,password,phone) "+
		"VALUES($1,$2,$3,$4) returning user_id;",
		name, email, getMD5Hash(password), phone).Scan(&userID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	u.logger.Info().Msgf("successfully added a user with ID: %d", userID)

	return nil
}

func (u *userEngine) UserEdit(userID int, name string, email string, password string, phone string) error {
	updateBusinessUserSQL := `
	UPDATE business_user
	SET name = $1, email = $2, password = $3, phone = $4
	WHERE user_id = $5;`

	_, err := u.sql.Exec(updateBusinessUserSQL, name, email, getMD5Hash(password), phone, userID)
	if err != nil {
		return err
	}

	return nil
}

func (u *userEngine) UserGet(userID int) (shared.BusinessUser, error) {
	rows, err := u.sql.Query("SELECT user_id, name, email, phone FROM business_user where "+
		"user_id = $1;", userID)
	if err != nil {
		return shared.BusinessUser{}, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var userInfo shared.BusinessUser
	if rows.Next() {
		err = rows.Scan(&userInfo.UserID, &userInfo.Name, &userInfo.Email, &userInfo.Phone)
		if err != nil {
			return shared.BusinessUser{}, helper.DatabaseError{DBError: err.Error()}
		}
	} else {
		return shared.BusinessUser{}, helper.UserError{Message: fmt.Sprintf("user with id: %d not found", userID)}
	}

	if err = rows.Err(); err != nil {
		return shared.BusinessUser{}, helper.DatabaseError{DBError: err.Error()}
	}

	return userInfo, nil
}

func (u *userEngine) CheckEmail(email string) (int, error) {
	rows, err := u.sql.Query("SELECT user_id FROM business_user where "+
		"email = $1;", email)
	if err != nil {
		return -1, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var userID int
	if rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return -1, helper.DatabaseError{DBError: err.Error()}
		}
	} else {
		// user does not exist for email
		return 0, nil
	}

	if err = rows.Err(); err != nil {
		return -1, helper.DatabaseError{DBError: err.Error()}
	}

	return userID, nil
}

func (u *userEngine) UserVerify(email string, password string) (shared.BusinessUser, error) {
	rows, err := u.sql.Query("SELECT user_id, name, email, phone FROM business_user where "+
		"email = $1 AND password = $2;", email, getMD5Hash(password))
	if err != nil {
		return shared.BusinessUser{}, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var userInfo shared.BusinessUser
	if rows.Next() {
		err = rows.Scan(&userInfo.UserID, &userInfo.Name, &userInfo.Email, &userInfo.Phone)
		if err != nil {
			return shared.BusinessUser{}, helper.DatabaseError{DBError: err.Error()}
		}
	} else {
		return shared.BusinessUser{}, helper.UserError{Message: fmt.Sprintf("user %s password mismatch", email)}
	}

	if err = rows.Err(); err != nil {
		return shared.BusinessUser{}, helper.DatabaseError{DBError: err.Error()}
	}

	return userInfo, nil
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (u *userEngine) ContactUs(phoneID string, uname string, email string, subject string, comments string) error {
	var contactUsID int
	err := u.sql.QueryRow("INSERT INTO contact_us(phone_id,uname,email,subject,comments,contact_date) "+
		"VALUES($1,$2,$3,$4,$5,$6) returning contact_us_id;",
		phoneID, uname, email, subject, comments, time.Now()).Scan(&contactUsID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	u.logger.Info().Msgf("successfully added a user comment with ID: %d", contactUsID)

	return nil
}

func (u *userEngine) ReportInaccurate(phoneID string, listingID int) error {
	var reportInaccurateId int
	err := u.sql.QueryRow("INSERT INTO report_inaccurate(phone_id,listing_id,report_inaccurate_add_date) "+
		"VALUES($1,$2,$3) returning report_inaccurate_id;",
		phoneID, listingID, time.Now()).Scan(&reportInaccurateId)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	u.logger.Info().Msgf("successfully registered ReportInaccurate with ID: %d", reportInaccurateId)
	return nil
}
