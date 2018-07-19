package user

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"

	"fmt"

	"github.com/phassans/banana/helper"
	"github.com/rs/xlog"
)

type userEngine struct {
	sql    *sql.DB
	logger xlog.Logger
}

type UserEngine interface {
	AddUser(name string, email string, password string, phone string) error
	VerifyUser(email string, password string) (int, error)
}

func NewUserEngine(psql *sql.DB, logger xlog.Logger) UserEngine {
	return &userEngine{psql, logger}
}

func (u *userEngine) AddUser(name string, email string, password string, phone string) error {
	// check business name unique
	userID, err := u.CheckEmail(email)
	if err != nil {
		return err
	}

	if userID != 0 {
		return helper.DuplicateEntity{Name: email}
	}

	err = u.sql.QueryRow("INSERT INTO business_user(name,email,password,phone) "+
		"VALUES($1,$2,$3,$4) returning user_id;",
		name, email, GetMD5Hash(password), phone).Scan(&userID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	u.logger.Infof("successfully added a user with ID: %d", userID)

	return nil
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
		err := rows.Scan(&userID)
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

func (u *userEngine) VerifyUser(email string, password string) (int, error) {
	rows, err := u.sql.Query("SELECT user_id FROM business_user where "+
		"email = $1 AND password = $2;", email, GetMD5Hash(password))
	if err != nil {
		return -1, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var userID int
	if rows.Next() {
		err := rows.Scan(&userID)
		if err != nil {
			return -1, helper.DatabaseError{DBError: err.Error()}
		}
	} else {
		return -1, helper.UserError{Message: fmt.Sprintf("user %s password mismatch", email)}
	}

	if err = rows.Err(); err != nil {
		return -1, helper.DatabaseError{DBError: err.Error()}
	}

	return userID, nil
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
