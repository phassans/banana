package model

type GenericEngine struct {
	BusinessEngine
	ListingEngine
	UserEngine
	FavoriteEngine
}

func NewGenericEngine(businessEngine BusinessEngine, userEngine UserEngine, listingEngine ListingEngine, favouriteEngine FavoriteEngine) Engine {
	return &GenericEngine{businessEngine, listingEngine, userEngine, favouriteEngine}
}

type Engine interface {
	BusinessEngine
	ListingEngine
	UserEngine
	FavoriteEngine
}
