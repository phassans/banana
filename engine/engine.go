package engine

type GenericEngine struct {
	BusinessEngine
	ListingEngine
	OwnerEngine
}

func NewGenericEngine(businessEngine BusinessEngine, ownerEngine OwnerEngine, listingEngine ListingEngine) Engine {
	return &GenericEngine{businessEngine, listingEngine, ownerEngine}
}

type Engine interface {
	BusinessEngine
	ListingEngine
	OwnerEngine
}
