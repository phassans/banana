package model

type GenericEngine struct {
	BusinessEngine
	ListingEngine
	UserEngine
}

func NewGenericEngine(businessEngine BusinessEngine, userEngine UserEngine, listingEngine ListingEngine) Engine {
	return &GenericEngine{businessEngine, listingEngine, userEngine}
}

type Engine interface {
	BusinessEngine
	ListingEngine
	UserEngine
}
