package model

type GenericEngine struct {
	BusinessEngine
	ListingEngine
	UserEngine
	FavoriteEngine
	NotificationEngine
}

func NewGenericEngine(businessEngine BusinessEngine, userEngine UserEngine, listingEngine ListingEngine,
	favouriteEngine FavoriteEngine, notificationEngine NotificationEngine) Engine {
	return &GenericEngine{businessEngine, listingEngine, userEngine, favouriteEngine, notificationEngine}
}

type Engine interface {
	BusinessEngine
	ListingEngine
	UserEngine
	FavoriteEngine
	NotificationEngine
}
