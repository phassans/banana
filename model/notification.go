package model

import (
	"database/sql"

	"github.com/phassans/banana/helper"
	"github.com/rs/xlog"
)

type (
	Notification struct {
		NotificationID     NotificationID
		PhoneId            string
		BusinessId         int
		Price              string
		Keywords           string
		DietaryRestriction []string
		Latitude           float64
		Longitude          float64
		Location           string
	}

	notificationEngine struct {
		sql            *sql.DB
		logger         xlog.Logger
		businessEngine BusinessEngine
	}

	NotificationEngine interface {
		AddNotification(notification Notification) error
		DeleteNotification(notificationID int) error
		GetAllNotifications(phoneID string) ([]Notification, error)
	}
)

func NewNotificationEngine(psql *sql.DB, logger xlog.Logger, businessEngine BusinessEngine) NotificationEngine {
	return &notificationEngine{psql, logger, businessEngine}
}

func (n *notificationEngine) AddNotification(notification Notification) error {
	var notificationID NotificationID

	// Add Notification
	err := n.sql.QueryRow("INSERT INTO notifications(phone_id,business_id,price,keywords) "+
		"VALUES($1,$2,$3,$4) returning notification_id;",
		notification.PhoneId, notification.BusinessId, notification.Price, notification.Keywords).Scan(&notificationID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}
	notification.NotificationID = notificationID

	// Add Notification Location
	err = n.AddNotificationLocation(notification)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	// Add Notification Dietary Restriction
	if len(notification.DietaryRestriction) > 0 {
		for _, restriction := range notification.DietaryRestriction {
			if err := n.AddNotificationDietaryRestriction(restriction, notification.NotificationID); err != nil {
				return err
			}
		}
	}

	n.logger.Infof("successfully added a notification with ID: %d", notificationID)
	return nil
}

func (n *notificationEngine) AddNotificationLocation(notification Notification) error {
	_, err := n.sql.Query("INSERT INTO notifications_location(notification_id,location,latitude,longitude) "+
		"VALUES($1,$2,$3,$4);",
		notification.NotificationID, notification.Location, notification.Latitude, notification.Longitude)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	n.logger.Infof("successfully added a notifications_location with ID: %d", notification.NotificationID)
	return nil
}

func (l *notificationEngine) AddNotificationDietaryRestriction(restriction string, notificationID NotificationID) error {
	addNotificationDietRestrictionSQL := "INSERT INTO notifications_dietary_restrictions(notification_id,restriction) " +
		"VALUES($1,$2);"

	_, err := l.sql.Query(addNotificationDietRestrictionSQL, notificationID, restriction)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	l.logger.Infof("add notifications_dietary_restrictions successful for notification:%d", notificationID)
	return nil
}

func (n *notificationEngine) DeleteNotification(notificationID int) error {
	if err := n.DeleteNotificationDietaryRestriction(notificationID); err != nil {
		return err
	}

	if err := n.DeleteNotificationLocation(notificationID); err != nil {
		return err
	}

	if err := n.DeleteNotificationInfo(notificationID); err != nil {
		return err
	}

	return nil
}

func (n *notificationEngine) DeleteNotificationInfo(notificationID int) error {
	sqlStatement := `DELETE FROM notifications WHERE notification_id = $1;`
	n.logger.Infof("deleting notification with query: %s and listing: %d", sqlStatement, notificationID)

	_, err := n.sql.Exec(sqlStatement, notificationID)
	return err
}

func (n *notificationEngine) DeleteNotificationLocation(notificationID int) error {
	sqlStatement := `DELETE FROM notifications_location WHERE notification_id = $1;`
	n.logger.Infof("deleting notifications_location with query: %s and listing: %d", sqlStatement, notificationID)

	_, err := n.sql.Exec(sqlStatement, notificationID)
	return err
}

func (n *notificationEngine) DeleteNotificationDietaryRestriction(notificationID int) error {
	sqlStatement := `DELETE FROM notifications_dietary_restrictions WHERE notification_id = $1;`
	n.logger.Infof("deleting notifications_location with query: %s and listing: %d", sqlStatement, notificationID)

	_, err := n.sql.Exec(sqlStatement, notificationID)
	return err
}

func (n *notificationEngine) GetAllNotifications(phoneID string) ([]Notification, error) {
	return nil, nil
}
