package notification

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"fmt"

	"github.com/phassans/banana/clients"
	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/model/business"
	"github.com/phassans/banana/shared"
	"github.com/rs/zerolog"
)

type (
	notificationEngine struct {
		sql            *sql.DB
		logger         zerolog.Logger
		businessEngine business.BusinessEngine
	}

	// NotificationEngine interface
	NotificationEngine interface {
		AddNotification(
			notificationName string,
			phoneID string,
			latitude float64,
			longitude float64,
			Location string,
			priceFilter string,
			dietaryFilter []string,
			distanceFilter string,
			keywords string,
		) error
		DeleteNotification(notificationID int) error
		GetAllNotifications(phoneID string) ([]shared.Notification, error)
		RegisterPhone(registrationToken string, phoneID string, phoneModel string) error
		GetStats() (string, error)
	}
)

// NewNotificationEngine returns an instance of notificationEngine
func NewNotificationEngine(psql *sql.DB, logger zerolog.Logger, businessEngine business.BusinessEngine) NotificationEngine {
	return &notificationEngine{psql, logger, businessEngine}
}

func (l *notificationEngine) isPhoneRegistered(phoneID string) (bool, error) {
	rows := l.sql.QueryRow("SELECT registration_token FROM register_phone where phone_id = $1;", phoneID)

	var registrationToken string
	err := rows.Scan(&registrationToken)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if registrationToken != "" {
		return true, nil
	}

	return false, nil
}

func (n *notificationEngine) RegisterPhone(registrationToken string, phoneID string, phoneModel string) error {
	registered, err := n.isPhoneRegistered(phoneID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	if registered {
		err := n.UpdatePhoneRegistrationToken(registrationToken, phoneID)
		if err != nil {
			return helper.DatabaseError{DBError: err.Error()}
		}
	} else {
		registerPhoneSQL := "INSERT INTO register_phone(registration_token,phone_id,phone_model,register_date) " +
			"VALUES($1,$2,$3,$4);"

		rows, err := n.sql.Query(registerPhoneSQL, registrationToken, phoneID, phoneModel, time.Now())
		if err != nil {
			return helper.DatabaseError{DBError: err.Error()}
		}
		defer rows.Close()

		n.logger.Info().Msgf("registerPhone successful for phoneID:%s", phoneID)
	}

	return nil
}

func (n *notificationEngine) GetStats() (string, error) {
	stats := make(map[string]int)

	iosStats, err := n.GetStat("ios")
	if err != nil {
		return "", err
	}
	stats["ios"] = iosStats

	androidStats, err := n.GetStat("android")
	if err != nil {
		return "", err
	}
	stats["android"] = androidStats

	jsonString, err := json.Marshal(stats)
	if err != nil {
		return "", err
	}

	return string(jsonString), nil
}

func (n *notificationEngine) GetStat(phoneModel string) (int, error) {
	var count int
	row := n.sql.QueryRow("SELECT COUNT(*) FROM register_phone where phone_model=$1", phoneModel)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return 0, nil
}

func (n *notificationEngine) UpdatePhoneRegistrationToken(registrationToken string, phoneID string) error {
	updatePhoneRegistrationToken := "UPDATE register_phone SET registration_token=$1, update_date=$2 WHERE phone_id=$3;"

	_, err := n.sql.Exec(updatePhoneRegistrationToken, registrationToken, time.Now(), phoneID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	n.logger.Info().Msgf("phone registered successfully: %s", phoneID)
	return nil
}

func (n *notificationEngine) AddNotification(
	notificationName string,
	phoneID string,
	latitude float64,
	longitude float64,
	location string,
	priceFilter string,
	dietaryFilter []string,
	distanceFilter string,
	keywords string,
) error {
	if location != "" {
		// getLatLonFromLocation
		resp, err := clients.GetLatLong(location)
		if err != nil {
			return err
		}
		latitude = resp.Lat
		longitude = resp.Lon
	}
	n.logger.Info().Msgf("notification location - lat:%f, lon:%f", latitude, longitude)

	var notificationID int
	// Add Notification
	err := n.sql.QueryRow("INSERT INTO notifications(notification_name,phone_id,latitude,longitude,location,price_filter,distance_filter,keywords) "+
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8) returning notification_id;",
		notificationName, phoneID, latitude, longitude, location, priceFilter, distanceFilter, keywords).Scan(&notificationID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	// Add Notification Dietary Restriction
	if len(dietaryFilter) > 0 {
		for _, restriction := range dietaryFilter {
			if err := n.AddNotificationDietaryRestriction(notificationID, restriction); err != nil {
				return err
			}
		}
	}

	n.logger.Info().Msgf("successfully added a notification with ID: %d", notificationID)
	return nil
}

func (n *notificationEngine) AddNotificationDietaryRestriction(notificationID int, restriction string) error {
	addNotificationDietRestrictionSQL := "INSERT INTO notifications_dietary_restrictions(notification_id,restriction) " +
		"VALUES($1,$2);"

	rows, err := n.sql.Query(addNotificationDietRestrictionSQL, notificationID, restriction)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	n.logger.Info().Msgf("add notifications_dietary_restrictions successful for notification:%d", notificationID)
	return nil
}

func (n *notificationEngine) DeleteNotification(notificationID int) error {
	if err := n.DeleteNotificationDietaryRestriction(notificationID); err != nil {
		return err
	}

	if err := n.DeleteNotificationInfo(notificationID); err != nil {
		return err
	}

	return nil
}

func (n *notificationEngine) DeleteNotificationInfo(notificationID int) error {
	sqlStatement := `DELETE FROM notifications WHERE notification_id = $1;`
	n.logger.Info().Msgf("deleting notification with query: %s and listing: %d", sqlStatement, notificationID)

	_, err := n.sql.Exec(sqlStatement, notificationID)
	return err
}

func (n *notificationEngine) DeleteNotificationDietaryRestriction(notificationID int) error {
	sqlStatement := `DELETE FROM notifications_dietary_restrictions WHERE notification_id = $1;`
	n.logger.Info().Msgf("deleting notifications_location with query: %s and listing: %d", sqlStatement, notificationID)

	_, err := n.sql.Exec(sqlStatement, notificationID)
	return err
}

func (n *notificationEngine) GetAllNotifications(phoneID string) ([]shared.Notification, error) {

	q := fmt.Sprintf("SELECT notification_id, notification_name, phone_id, latitude, longitude, location, price_filter, distance_filter, keywords FROM notifications where phone_id = '%s';", phoneID)
	rows, err := n.sql.Query(q)
	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	i := 1
	var notifications []shared.Notification
	for rows.Next() {
		var notification shared.Notification
		err = rows.Scan(
			&notification.NotificationID,
			&notification.NotificationName,
			&notification.PhoneID,
			&notification.Latitude,
			&notification.Longitude,
			&notification.Location,
			&notification.PriceFilter,
			&notification.DistanceFilter,
			&notification.Keywords,
		)
		if err != nil {
			return nil, helper.DatabaseError{DBError: err.Error()}
		}
		if notification.NotificationName == "" {
			notification.NotificationName = fmt.Sprintf("Notification %d", i)
			i++
		}
		notifications = append(notifications, notification)
	}

	if err = rows.Err(); err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	// addDietaryRestrictionsToListings
	notifications, err = n.AddDietaryRestrictionsToNotifications(notifications)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (n *notificationEngine) AddDietaryRestrictionsToNotifications(notifications []shared.Notification) ([]shared.Notification, error) {
	// get dietary restriction
	for i := 0; i < len(notifications); i++ {
		// add dietary restriction
		rests, err := n.GetNotificationsDietaryRestriction(notifications[i].NotificationID)
		if err != nil {
			return nil, err
		}
		notifications[i].DietaryFilters = rests
	}
	return notifications, nil
}

func (n *notificationEngine) GetNotificationsDietaryRestriction(notificationID int) ([]string, error) {
	rows, err := n.sql.Query("SELECT restriction FROM notifications_dietary_restrictions WHERE notification_id = $1", notificationID)
	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var rests []string
	for rows.Next() {
		var rest string
		err = rows.Scan(&rest)
		if err != nil {
			return nil, helper.DatabaseError{DBError: err.Error()}
		}
		rests = append(rests, rest)
	}

	if err = rows.Err(); err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	return rests, nil
}
