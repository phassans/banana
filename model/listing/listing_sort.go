package listing

import (
	"fmt"
	"strings"
	"time"

	"github.com/phassans/banana/shared"
	"github.com/umahmood/haversine"
)

type (
	sortDistanceView struct {
		listing shared.Listing
		mile    float64
	}

	sortPriceView struct {
		listing shared.Listing
		price   float64
	}

	sortTimeView struct {
		listing  shared.Listing
		timeLeft float64
	}

	CurrentLocation struct {
		Latitude  float64
		Longitude float64
	}
)

func (l *listingEngine) SortListings(listings []shared.Listing, sortingType string,
	currentLocation CurrentLocation) ([]shared.Listing, error) {

	if sortingType == "distance" || sortingType == "" {
		return l.SortListingsByDistance(listings, currentLocation)
	} else if sortingType == "price" {
		return l.SortListingsByPrice(listings)
	} else if sortingType == "timeLeft" {
		return l.SortListingsByTimeLeft(listings)
	}

	return nil, nil
}

func (l *listingEngine) SortListingsByTimeLeft(listings []shared.Listing) ([]shared.Listing, error) {
	var ll []sortTimeView
	for _, listing := range listings {

		dateTime := GetListingDateTime(listing.StartDate, listing.StartTime)
		then, err := time.Parse(shared.DateTimeFormat, dateTime)
		if err != nil {
			return nil, nil
		}

		duration := time.Since(then)

		s := sortTimeView{listing: listing, timeLeft: duration.Seconds()}
		ll = append(ll, s)
	}

	// sort
	priceView := l.OrderListingsByTime(ll)

	// put in listing struct
	var listingsResult []shared.Listing
	for _, view := range priceView {
		listingsResult = append(listingsResult, view.listing)
	}

	return listingsResult, nil
}

func GetListingDateTime(endDate string, endTime string) string {
	listingEndDate := strings.Split(endDate, "T")[0]
	listingEndTime := strings.Split(endTime, "T")[1]
	return fmt.Sprintf("%sT%s", listingEndDate, listingEndTime)
}

func (l *listingEngine) SortListingsByPrice(listings []shared.Listing) ([]shared.Listing, error) {
	var ll []sortPriceView
	for _, listing := range listings {
		s := sortPriceView{listing: listing, price: listing.NewPrice}
		ll = append(ll, s)
	}

	// sort
	priceView := l.OrderListingsByPrice(ll)

	// put in listing struct
	var listingsResult []shared.Listing
	for _, view := range priceView {
		listingsResult = append(listingsResult, view.listing)
	}

	return listingsResult, nil
}

func (l *listingEngine) SortListingsByDistance(listings []shared.Listing, currentLocation CurrentLocation) ([]shared.Listing, error) {
	var ll []sortDistanceView
	for _, listing := range listings {
		// get LatLon
		geo, err := l.GetListingsLatLon(listing.BusinessID)
		if err != nil {
			return nil, err
		}

		// append latLon
		fromMobile := haversine.Coord{Lat: currentLocation.Latitude, Lon: currentLocation.Longitude}
		fromDB := haversine.Coord{Lat: geo.Latitude, Lon: geo.Longitude}
		mi, _ := haversine.Distance(fromMobile, fromDB)
		listing.DistanceFromLocation = mi

		fmt.Printf("business_id: %d and distance: %f \n", listing.BusinessID, mi)

		s := sortDistanceView{listing: listing, mile: mi}
		ll = append(ll, s)
	}

	// sort
	distanceView := l.OrderListingsByDistance(ll)

	// put in listing struct
	var listingsResult []shared.Listing
	for _, view := range distanceView {
		listingsResult = append(listingsResult, view.listing)
	}

	return listingsResult, nil
}
