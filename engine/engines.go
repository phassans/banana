package engine

type ListingEngine interface {
	AddListing(title string, description string, price float64, startTime string, endTime string, businessName string) error
}

type BusinessEngine interface {
	AddBusiness(businessName string, phone string, website string) (int, error)
	AddBusinessAddress(line1 string, line2 string, city string, postalCode string, state string, country string, businessName string, otherDetails string) error
	AddGeoInfo(address string, addressID int, businessID int) error
	GetBusinessIDFromName(businessName string) (int, error)
}

type OwnerEngine interface {
	AddOwner(firstName string, lastName string, phone string, email string, businessName string) error
}

type Engine interface {
	BusinessEngine
	ListingEngine
	OwnerEngine
}
