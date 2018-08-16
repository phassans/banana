package notification

import (
	"database/sql"

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
	}
)

// NewNotificationEngine returns an instance of notificationEngine
func NewNotificationEngine(psql *sql.DB, logger zerolog.Logger, businessEngine business.BusinessEngine) NotificationEngine {
	return &notificationEngine{psql, logger, businessEngine}
}

func (n *notificationEngine) RegisterPhone(registrationToken string, phoneID string, phoneModel string) error {

	registerPhoneSQL := "INSERT INTO register_phone(registration_token,phone_id,phone_model) " +
		"VALUES($1,$2,$3);"

	rows, err := n.sql.Query(registerPhoneSQL, registrationToken, phoneID, phoneModel)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	n.logger.Info().Msgf("registerPhone successful for phoneID:%s", phoneID)
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

	return notifications, nil
}
