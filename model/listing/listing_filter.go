package listing

import (
	"strconv"

	"github.com/phassans/banana/shared"
)

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

func (l *listingEngine) FilterByDistance(listings []shared.Listing, distanceFilter string) ([]shared.Listing, error) {
	if distanceFilter == "all" {
		distanceFilter = "50"
	}

	// if parse error, do not apply filter, just return
	dFilter, err := strconv.ParseFloat(distanceFilter, 64)
	if err != nil {
		return listings, err
	}

	// get dietary restriction
	var listingsResult []shared.Listing
	for _, listing := range listings {
		if listing.DistanceFromLocation <= dFilter {
			listingsResult = append(listingsResult, listing)
		}
	}
	return listingsResult, nil
}

func (l *listingEngine) FilterByDietaryRestrictions(listings []shared.Listing, dietaryFilters []string) ([]shared.Listing, error) {
	// get dietary restriction
	var listingsResult []shared.Listing
	for _, listing := range listings {
		listingAdded := false
		rests, err := l.GetListingsDietaryRestriction(listing.ListingID)
		if err != nil {
			return nil, err
		}

		for _, rest := range rests {
			if listingAdded {
				continue
			}
			for _, dietaryFilter := range dietaryFilters {
				if rest == dietaryFilter {
					listingsResult = append(listingsResult, listing)
					listingAdded = true
					continue
				}
			}
		}
	}
	return listingsResult, nil
}

func (f *listingEngine) getListingStatus(listing shared.Listing) string {

	startDateTimeLeft, err := CalculateTimeLeft(listing.StartDate, listing.EndTime)
	if err != nil {
		return ""
	}

	if startDateTimeLeft > 1 {
		return "scheduled"
	} else if startDateTimeLeft < 0 && !listing.MultipleDays && !listing.Recurring {
		return "ended"
	}

	if listing.MultipleDays {
		endDateTimeLeft, err := CalculateTimeLeft(listing.EndDate, listing.EndTime)
		if err != nil {
			return ""
		}

		if endDateTimeLeft < 0 {
			return "ended"
		}
	} else if listing.Recurring {
		recurringEndDateTimeLeft, err := CalculateTimeLeft(listing.RecurringEndDate, listing.EndTime)
		if err != nil {
			return ""
		}

		if recurringEndDateTimeLeft < 0 {
			return "ended"
		}
	}

	return "active"
}

func (f *listingEngine) filterListingBasedOnStatus(listings []shared.Listing, status string) []shared.Listing {
	if status == "all" || status == "" {
		return listings
	}

	var resultListings []shared.Listing
	for _, listing := range listings {
		if listing.ListingStatus == status {
			resultListings = append(resultListings, listing)
		}
	}

	return resultListings
}
