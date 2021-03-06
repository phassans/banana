package listing

import (
	"fmt"
	"sort"
	"time"

	"strings"

	"database/sql"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/model/common"
	"github.com/phassans/banana/shared"
	"github.com/umahmood/haversine"
)

type (
	sortListingEngine struct {
		listings        []shared.Listing
		sortingType     string
		currentLocation shared.GeoLocation
		sql             *sql.DB
	}

	// SortListingEngine interface
	SortListingEngine interface {
		SortListings(bool, string, bool, bool) ([]shared.Listing, error)
	}
)

// NewSortListingEngine returns an instance of sortListingEngine
func NewSortListingEngine(listings []shared.Listing, sortingType string,
	currentLocation shared.GeoLocation, sql *sql.DB) SortListingEngine {
	return &sortListingEngine{listings, sortingType, currentLocation, sql}
}

func (l *sortListingEngine) SortListings(isFuture bool, searchDay string, isSearch bool, isFavorite bool) ([]shared.Listing, error) {

	// have to sort by distance, in order to calculate distanceFromLocation
	if l.currentLocation.Latitude != 0 && l.currentLocation.Longitude != 0 {
		l.sortListingsByDistance(isFuture, searchDay, isFavorite)
	}

	// for future listings always sort by timeLeft
	if l.sortingType == "" && isFuture {
		l.sortListingsByTimeLeft()
		return l.listings, nil
	}

	if l.sortingType == shared.SortByPrice {
		l.sortListingsByPrice()
	} else if l.sortingType == shared.SortByTimeLeft {
		l.sortListingsByTimeLeft()
	} else if l.sortingType == shared.SortByDateAdded {
		if isFavorite {
			l.sortListingsByFavoriteDateAdded()
		} else {
			l.sortListingsByDateAdded()
		}
	} else if l.sortingType == shared.SortByMostPopular {
		l.sortListingsByMostPopular()
	}

	return l.listings, nil
}

func (l *sortListingEngine) sortListingsByTimeLeft() error {
	var ll []shared.SortView
	for _, listing := range l.listings {

		timeLeft, err := calculateTimeLeft(listing.ListingDate, listing.EndTime, l.currentLocation)
		if err != nil {
			return err
		}

		s := shared.SortView{Listing: listing, TimeLeft: float64(timeLeft)}
		ll = append(ll, s)
	}

	// put in listing struct
	var listingsResult []shared.Listing
	for _, view := range l.orderListings(ll, shared.SortByTimeLeft) {
		listingsResult = append(listingsResult, view.Listing)
	}
	l.listings = listingsResult
	return nil
}

func (l *sortListingEngine) sortListingsByMostPopular() error {
	var ll []shared.SortView
	for _, listing := range l.listings {
		s := shared.SortView{Listing: listing, UpVotes: listing.UpVotes}
		ll = append(ll, s)
	}

	// put in listing struct
	var listingsResult []shared.Listing
	for _, view := range l.orderListings(ll, shared.SortByMostPopular) {
		listingsResult = append(listingsResult, view.Listing)
	}
	l.listings = listingsResult
	return nil
}

func (l *sortListingEngine) sortListingsByPrice() error {
	var ll []shared.SortView
	var discountedListings []shared.Listing

	for _, listing := range l.listings {
		if listing.NewPrice == 0 {
			discountedListings = append(discountedListings, listing)
			continue
		}

		s := shared.SortView{Listing: listing, Price: listing.NewPrice}
		ll = append(ll, s)
	}

	// put in listing struct
	var listingsResult []shared.Listing
	for _, view := range l.orderListings(ll, shared.SortByPrice) {
		listingsResult = append(listingsResult, view.Listing)
	}

	// append two lists
	listingsResult = append(listingsResult, discountedListings...)

	l.listings = listingsResult
	return nil
}

func (l *sortListingEngine) sortListingsByDistance(isFuture bool, searchDay string, isFavorite bool) error {
	var ll []shared.SortView

	for _, listing := range l.listings {
		// append latLon
		fromMobile := haversine.Coord{Lat: l.currentLocation.Latitude, Lon: l.currentLocation.Longitude}
		fromDB := haversine.Coord{Lat: listing.Latitude, Lon: listing.Longitude}
		mi, _ := haversine.Distance(fromMobile, fromDB)
		if isFavorite || mi <= getMaxDistance(isFuture) {
			listing.DistanceFromLocation = mi
			s := shared.SortView{Listing: listing, Mile: mi}
			ll = append(ll, s)
		}
	}

	// if favorites do not sort ByDistance
	if l.sortingType != "" && l.sortingType != shared.SortByDistance {
		var listingsResult []shared.Listing
		for _, view := range ll {
			listingsResult = append(listingsResult, view.Listing)
		}
		l.listings = listingsResult
		return nil
	}

	// put in listing struct
	var listingsResult = make([]shared.Listing, 0)
	for _, view := range l.orderListings(ll, shared.SortByDistance) {
		listingsResult = append(listingsResult, view.Listing)
	}
	l.listings = listingsResult

	return l.GroupListingsBasedOnSearchDay(searchDay)
}

func (l *sortListingEngine) GroupListingsBasedOnSearchDay(searchDay string) error {
	// group them by day
	if searchDay != "" {
		var days []string

		timeInZone, err := getCurrentTimeInTimeZone(l.currentLocation)
		if err != nil {
			return err
		}

		currentDate := timeInZone.Format(shared.DateFormat)
		curDate, err := time.Parse(shared.DateFormat, currentDate)
		if err != nil {
			return err
		}

		startDay := 0
		endDay := 0
		switch searchDay {
		case shared.SearchThisWeek:
			//fmt.Println("listingDate.Weekday().String()", curDate.Weekday().String())
			endDay = 6 - shared.DayMap[strings.ToLower(curDate.Weekday().String())]
		case shared.SearchNextWeek:
			startDay = 6 - shared.DayMap[strings.ToLower(curDate.Weekday().String())]
			endDay = 13 - shared.DayMap[strings.ToLower(curDate.Weekday().String())]
		default:
			return nil
		}

		if searchDay == shared.SearchThisWeek {
			days = append(days, curDate.Weekday().String())
			for i := 0; i < endDay; i++ {
				nextDate := curDate.Add(time.Hour * 24)
				days = append(days, nextDate.Weekday().String())
				curDate = nextDate
			}
		} else if searchDay == shared.SearchNextWeek {
			for i := 0; i < endDay; i++ {
				nextDate := curDate.Add(time.Hour * 24)
				if i >= startDay {
					days = append(days, nextDate.Weekday().String())
				}
				curDate = nextDate
			}
		}

		var fResult []shared.Listing
		for _, day := range days {
			for _, listing := range l.listings {
				listingDateFormatted, err := time.Parse(shared.DateFormatSQL, strings.Split(listing.ListingDate, "T")[0])
				if err != nil {
					return nil
				}
				if listingDateFormatted.Weekday().String() == day {
					fResult = append(fResult, listing)
				}
			}
		}
		l.listings = fResult
	}

	return nil
}

func getMaxDistance(isFuture bool) float64 {
	if isFuture {
		return common.MaxDistanceForFutureDeals
	}
	return common.MaxDistanceForTodaysDeals
}

func (l *sortListingEngine) sortListingsByDateAdded() error {
	var ll []shared.SortView
	for _, listing := range l.listings {
		listingDateFormatted, err := time.Parse(shared.DateTimeFormat, listing.ListingCreateDate)
		if err != nil {
			return err
		}
		s := shared.SortView{Listing: listing, ListingDate: listingDateFormatted}
		ll = append(ll, s)
	}

	// put in listing struct
	var listingsResult []shared.Listing
	for _, view := range l.orderListings(ll, shared.SortByDateAdded) {
		listingsResult = append(listingsResult, view.Listing)
	}

	return nil
}

func (l *sortListingEngine) sortListingsByFavoriteDateAdded() error {
	var ll []shared.SortView
	for _, listing := range l.listings {
		listingDateFormatted, err := time.Parse(shared.DateTimeFormat, listing.Favorite.FavoriteAddDate)
		if err != nil {
			return err
		}
		s := shared.SortView{Listing: listing, ListingDate: listingDateFormatted}
		ll = append(ll, s)
	}

	// put in listing struct
	var listingsResult []shared.Listing
	for _, view := range l.orderListings(ll, shared.SortByDateAdded) {
		listingsResult = append(listingsResult, view.Listing)
	}

	return nil
}

func (l *sortListingEngine) orderListings(listings []shared.SortView, orderType string) []shared.SortView {
	switch orderType {
	case shared.SortByTimeLeft:
		sort.Slice(listings, func(i, j int) bool {
			return listings[i].TimeLeft < listings[j].TimeLeft
		})
		return listings
	case shared.SortByPrice:
		sort.Slice(listings, func(i, j int) bool {
			return listings[i].Price < listings[j].Price
		})
		return listings
	case shared.SortByDistance:
		sort.Slice(listings, func(i, j int) bool {
			return listings[i].Mile < listings[j].Mile
		})
		return listings
	case shared.SortByDateAdded:
		sort.Slice(listings, func(i, j int) bool {
			//fmt.Printf("%t\n", listings[j].ListingDate.Before(listings[i].ListingDate))
			return listings[j].ListingDate.Before(listings[i].ListingDate)
		})
		return listings
	case shared.SortByMostPopular:
		sort.Slice(listings, func(i, j int) bool {
			return listings[i].UpVotes > listings[j].UpVotes
		})
		return listings
	}
	return nil
}

func (l *sortListingEngine) GetAllListingsLatLon(businessIDs []string) ([]shared.AddressGeo, error) {

	businesses := strings.Join(businessIDs, ",")

	businessesQuery := fmt.Sprintf("SELECT address_id, business_id, latitude, longitude  FROM business_address WHERE business_id IN (%s)", businesses)
	//fmt.Println(businessesQuery)

	rows, err := l.sql.Query(businessesQuery)
	defer rows.Close()

	if err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	var geos []shared.AddressGeo
	for rows.Next() {
		geo := shared.AddressGeo{}
		err = rows.Scan(&geo.AddressID, &geo.BusinessID, &geo.Latitude, &geo.Longitude)
		if err != nil {
			return nil, helper.DatabaseError{DBError: err.Error()}
		}
		geos = append(geos, geo)
	}

	if err = rows.Err(); err != nil {
		return nil, helper.DatabaseError{DBError: err.Error()}
	}

	return geos, nil
}

func (l *sortListingEngine) GetListingsLatLon(businessID int) (shared.AddressGeo, error) {
	rows, err := l.sql.Query("SELECT address_id, business_id, latitude, longitude  FROM business_address WHERE business_id = $1", businessID)
	defer rows.Close()

	if err != nil {
		return shared.AddressGeo{}, helper.DatabaseError{DBError: err.Error()}
	}

	defer rows.Close()

	geo := shared.AddressGeo{}
	if rows.Next() {
		err = rows.Scan(&geo.AddressID, &geo.BusinessID, &geo.Latitude, &geo.Longitude)
		if err != nil {
			return shared.AddressGeo{}, helper.DatabaseError{DBError: err.Error()}
		}
	}

	if err = rows.Err(); err != nil {
		return shared.AddressGeo{}, helper.DatabaseError{DBError: err.Error()}
	}

	return geo, nil
}

func getListingDateTime(endDate string, endTime string) string {
	listingEndDate := strings.Split(endDate, "T")[0]
	listingEndTime := strings.Split(endTime, "T")[1]
	return fmt.Sprintf("%sT%s", listingEndDate, listingEndTime)
}
