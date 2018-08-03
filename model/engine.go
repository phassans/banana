package model

import (
	"github.com/phassans/banana/model/business"
	"github.com/phassans/banana/model/favourite"
	"github.com/phassans/banana/model/listing"
	"github.com/phassans/banana/model/notification"
	"github.com/phassans/banana/model/user"
)

type genericEngine struct {
	business.BusinessEngine
	listing.ListingEngine
	user.UserEngine
	favourite.FavoriteEngine
	notification.NotificationEngine
}

// NewGenericEngine returns genericEngine
func NewGenericEngine(businessEngine business.BusinessEngine, userEngine user.UserEngine, listingEngine listing.ListingEngine,
	favouriteEngine favourite.FavoriteEngine, notificationEngine notification.NotificationEngine) Engine {
	return &genericEngine{businessEngine, listingEngine, userEngine, favouriteEngine, notificationEngine}
}

// Engine common engine interface
type Engine interface {
	business.BusinessEngine
	listing.ListingEngine
	user.UserEngine
	favourite.FavoriteEngine
	notification.NotificationEngine
}
