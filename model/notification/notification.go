package notification

import (
	"database/sql"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/model/business"
	"github.com/phassans/banana/shared"
	"github.com/rs/xlog"
)

type (
	notificationEngine struct {
		sql            *sql.DB
		logger         xlog.Logger
		businessEngine business.BusinessEngine
	}

	NotificationEngine interface {
		AddNotification(notification shared.Notification) error
		DeleteNotification(notificationID int) error
		GetAllNotifications(phoneID string) ([]shared.Notification, error)
	}
)

func NewNotificationEngine(psql *sql.DB, logger xlog.Logger, businessEngine business.BusinessEngine) NotificationEngine {
	return &notificationEngine{psql, logger, businessEngine}
}

func (n *notificationEngine) AddNotification(notification shared.Notification) error {
	var notificationID int

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

func (n *notificationEngine) AddNotificationLocation(notification shared.Notification) error {
	_, err := n.sql.Query("INSERT INTO notifications_location(notification_id,location,latitude,longitude) "+
		"VALUES($1,$2,$3,$4);",
		notification.NotificationID, notification.Location, notification.Latitude, notification.Longitude)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}

	n.logger.Infof("successfully added a notifications_location with ID: %d", notification.NotificationID)
	return nil
}

func (l *notificationEngine) AddNotificationDietaryRestriction(restriction string, notificationID int) error {
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

func (n *notificationEngine) GetAllNotifications(phoneID string) ([]shared.Notification, error) {
	return nil, nil
}