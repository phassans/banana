package listing

import "github.com/phassans/banana/shared"

func (l *listingEngine) FilterByPrice(listings []shared.Listing, priceFilter float64) ([]shared.Listing, error) {
	// get dietary restriction
	var listingsResult []shared.Listing
	for _, listing := range listings {
		if listing.NewPrice <= priceFilter {
			listingsResult = append(listingsResult, listing)
		}
	}
	return listingsResult, nil
}

func (l *listingEngine) FilterByDietaryRestrictions(listings []shared.Listing, dietaryFilter string) ([]shared.Listing, error) {
	// get dietary restriction
	var listingsResult []shared.Listing
	for _, listing := range listings {
		rests, err := l.GetListingsDietaryRestriction(listing.ListingID)
		if err != nil {
			return nil, err
		}
		for _, rest := range rests {
			if rest == dietaryFilter {
				listingsResult = append(listingsResult, listing)
			}
		}
	}
	return listingsResult, nil
}
