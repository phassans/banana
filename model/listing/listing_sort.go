package listing

import (
	"fmt"
	"sort"

	"strings"

	"database/sql"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
	"github.com/umahmood/haversine"
)

type (
	sortListingEngine struct {
		listings        []shared.Listing
		sortingType     string
		currentLocation shared.CurrentLocation
		sql             *sql.DB
	}

	SortListingEngine interface {
		SortListings() ([]shared.Listing, error)
	}
)

const (
	sortByDistance = "distance"
	sortByPrice    = "price"
	sortByTimeLeft = "timeLeft"
)

func NewSortListingEngine(listings []shared.Listing, sortingType string,
	currentLocation shared.CurrentLocation, sql *sql.DB) SortListingEngine {
	return &sortListingEngine{listings, sortingType, currentLocation, sql}
}

func (l *sortListingEngine) SortListings() ([]shared.Listing, error) {

	// have to sort by distance, in order to calculate distanceFromLocation
	l.SortListingsByDistance()

	if l.sortingType == sortByPrice {
		l.SortListingsByPrice()
	} else if l.sortingType == sortByTimeLeft {
		l.SortListingsByTimeLeft()
	}

	return l.listings, nil
}

func (l *sortListingEngine) SortListingsByTimeLeft() error {
	var ll []shared.SortView
	for _, listing := range l.listings {

		timeLeft, err := calculateTimeLeft(listing.ListingDate, listing.EndTime)
		if err != nil {
			return err
		}

		s := shared.SortView{Listing: listing, TimeLeft: float64(timeLeft)}
		ll = append(ll, s)
	}

	// put in listing struct
	var listingsResult []shared.Listing
	for _, view := range l.orderListings(ll, sortByTimeLeft) {
		listingsResult = append(listingsResult, view.Listing)
	}
	l.listings = listingsResult
	return nil
}

func (l *sortListingEngine) SortListingsByPrice() error {
	var ll []shared.SortView
	for _, listing := range l.listings {
		s := shared.SortView{Listing: listing, Price: listing.NewPrice}
		ll = append(ll, s)
	}

	// put in listing struct
	var listingsResult []shared.Listing
	for _, view := range l.orderListings(ll, sortByPrice) {
		listingsResult = append(listingsResult, view.Listing)
	}
	l.listings = listingsResult

	return nil
}

func (l *sortListingEngine) SortListingsByDistance() error {
	var ll []shared.SortView
	for _, listing := range l.listings {
		// get LatLon
		geo, err := l.GetListingsLatLon(listing.BusinessID)
		if err != nil {
			return err
		}

		// append latLon
		fromMobile := haversine.Coord{Lat: l.currentLocation.Latitude, Lon: l.currentLocation.Longitude}
		fromDB := haversine.Coord{Lat: geo.Latitude, Lon: geo.Longitude}
		mi, _ := haversine.Distance(fromMobile, fromDB)
		listing.DistanceFromLocation = mi

		s := shared.SortView{Listing: listing, Mile: mi}
		ll = append(ll, s)
	}

	// put in listing struct
	var listingsResult []shared.Listing
	for _, view := range l.orderListings(ll, sortByDistance) {
		listingsResult = append(listingsResult, view.Listing)
	}
	l.listings = listingsResult

	return nil
}

func (l *sortListingEngine) orderListings(listings []shared.SortView, orderType string) []shared.SortView {
	switch orderType {
	case sortByTimeLeft:
		sort.Slice(listings, func(i, j int) bool {
			return listings[i].TimeLeft < listings[j].TimeLeft
		})
		return listings
	case sortByPrice:
		sort.Slice(listings, func(i, j int) bool {
			return listings[i].Price < listings[j].Price
		})
		return listings
	case sortByDistance:
		sort.Slice(listings, func(i, j int) bool {
			return listings[i].Mile < listings[j].Mile
		})
		return listings
	}
	return nil
}

func (l *sortListingEngine) GetListingsLatLon(businessID int) (shared.AddressGeo, error) {
	rows, err := l.sql.Query("SELECT address_id, business_id, latitude, longitude  FROM address_geo WHERE business_id = $1", businessID)
	if err != nil {
		return shared.AddressGeo{}, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	geo := shared.AddressGeo{}
	if rows.Next() {
		err := rows.Scan(&geo.AddressID, &geo.BusinessID, &geo.Latitude, &geo.Longitude)
		if err != nil {
			return shared.AddressGeo{}, helper.DatabaseError{DBError: err.Error()}
		}
	}

	if err = rows.Err(); err != nil {
		return shared.AddressGeo{}, helper.DatabaseError{DBError: err.Error()}
	}

	return geo, nil
}

func GetListingDateTime(endDate string, endTime string) string {
	listingEndDate := strings.Split(endDate, "T")[0]
	listingEndTime := strings.Split(endTime, "T")[1]
	return fmt.Sprintf("%sT%s", listingEndDate, listingEndTime)
}
