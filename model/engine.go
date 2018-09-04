package model

import (
	"github.com/phassans/banana/model/business"
	"github.com/phassans/banana/model/favourite"
	"github.com/phassans/banana/model/listing"
	"github.com/phassans/banana/model/notification"
	"github.com/phassans/banana/model/prefernce"
	"github.com/phassans/banana/model/user"
)

type genericEngine struct {
	business.BusinessEngine
	listing.ListingEngine
	user.UserEngine
	favourite.FavoriteEngine
	notification.NotificationEngine
	prefernce.PreferenceEngine
}

// NewGenericEngine returns genericEngine
func NewGenericEngine(businessEngine business.BusinessEngine, userEngine user.UserEngine, listingEngine listing.ListingEngine,
	favouriteEngine favourite.FavoriteEngine, notificationEngine notification.NotificationEngine, preferenceEngine prefernce.PreferenceEngine) Engine {
	return &genericEngine{businessEngine, listingEngine, userEngine, favouriteEngine, notificationEngine, preferenceEngine}
}

// Engine common engine interface
type Engine interface {
	business.BusinessEngine
	listing.ListingEngine
	user.UserEngine
	favourite.FavoriteEngine
	notification.NotificationEngine
	prefernce.PreferenceEngine
}
