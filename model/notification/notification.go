package notification

import (
	"database/sql"

	"github.com/phassans/banana/clients"
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
	if notification.Location != "" {
		// getLatLonFromLocation
		resp, err := clients.GetLatLong(notification.Location)
		if err != nil {
			return err
		}
		notification.Latitude = resp.Lat
		notification.Longitude = resp.Lon
	}
	xlog.Infof("notification location - lat:%f, lon:%f", notification.Latitude, notification.Longitude)

	var notificationID int

	// Add Notification
	err := n.sql.QueryRow("INSERT INTO notifications(phone_id,business_id,price,keywords,location,latitude,longitude) "+
		"VALUES($1,$2,$3,$4,$5,$6,$7) returning notification_id;",
		notification.PhoneId, notification.BusinessId, notification.Price, notification.Keywords,
		notification.Location, notification.Latitude, notification.Longitude).Scan(&notificationID)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}
	notification.NotificationID = notificationID

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

func (l *notificationEngine) AddNotificationDietaryRestriction(restriction string, notificationID int) error {
	addNotificationDietRestrictionSQL := "INSERT INTO notifications_dietary_restrictions(notification_id,restriction) " +
		"VALUES($1,$2);"

	rows, err := l.sql.Query(addNotificationDietRestrictionSQL, notificationID, restriction)
	if err != nil {
		return helper.DatabaseError{DBError: err.Error()}
	}
	defer rows.Close()

	l.logger.Infof("add notifications_dietary_restrictions successful for notification:%d", notificationID)
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
	n.logger.Infof("deleting notification with query: %s and listing: %d", sqlStatement, notificationID)

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
	rows, err := n.sql.Query("SELECT notification_id, phone_id, business_id, price, keywords, location, latitude, longitude "+
		"FROM notifications "+
		"where phone_id = $1;", phoneID)
	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var notifications []shared.Notification
	for rows.Next() {
		var notific shared.Notification
		err := rows.Scan(
			&notific.NotificationID,
			&notific.PhoneId,
			&notific.BusinessId,
			&notific.Price,
			&notific.Keywords,
			&notific.Location,
			&notific.Latitude,
			&notific.Longitude,
		)
		if err != nil {
			return nil, helper.DatabaseError{DBError: err.Error()}
		}
		notifications = append(notifications, notific)
	}

	if err = rows.Err(); err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	return notifications, nil
}
