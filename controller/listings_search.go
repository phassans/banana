package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
)

type (
	listingsSearchRequest struct {
		Future         bool     `json:"future"`
		Search         bool     `json:"search,omitempty"`
		ListingTypes   []string `json:"listingTypes,omitempty"`
		Latitude       float64  `json:"latitude,omitempty"`
		Longitude      float64  `json:"longitude,omitempty"`
		Location       string   `json:"location,omitempty"`
		PriceFilter    float64  `json:"priceFilter,omitempty"`
		DietaryFilters []string `json:"dietaryFilters,omitempty"`
		DistanceFilter string   `json:"distanceFilter,omitempty"`
		Keywords       string   `json:"keywords,omitempty"`
		SortBy         string   `json:"sortBy,omitempty"`
		SearchDay      string   `json:"searchDay,omitempty"`

		PhoneID string `json:"phoneId"`
	}

	listingsSearchResult struct {
		Result  []shared.SearchListingResult
		Message string    `json:"message,omitempty"`
		Error   *APIError `json:"error,omitempty"`
	}

	listingsSearchEndpoint struct{}
)

var listingsSearch postEndpoint = listingsSearchEndpoint{}

func (r listingsSearchEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(listingsSearchRequest)

	logger := shared.GetLogger()
	logger = logger.With().
		Str("endpoint", r.GetPath()).
		Bool("future", request.Future).
		Bool("search", request.Search).
		Strs("listingTypes", request.ListingTypes).
		Float64("latitude", request.Latitude).
		Float64("longitude", request.Longitude).
		Str("location", request.Location).
		Float64("priceFilter", request.PriceFilter).
		Strs("dietaryFilters", request.DietaryFilters).
		Str("distanceFilter", request.DistanceFilter).
		Str("keywords", request.Keywords).
		Str("searchDay", request.SearchDay).
		Str("sortBy", request.SortBy).Logger()
	logger.Info().Msgf("search request")

	if request.SortBy == "dateadded" {
		request.SortBy = shared.SortByDateAdded
	}

	// validate input
	if err := r.Validate(request); err != nil {
		return nil, err
	}

	// if search keyword is "all cuisines" then search for everything
	if strings.ToLower(request.Keywords) == "all cuisines" {
		request.Keywords = ""
	}

	searchRequest := shared.SearchRequest{
		ListingTypes:   request.ListingTypes,
		Future:         request.Future,
		Latitude:       request.Latitude,
		Longitude:      request.Longitude,
		Location:       request.Location,
		PriceFilter:    request.PriceFilter,
		DietaryFilters: request.DietaryFilters,
		DistanceFilter: request.DistanceFilter,
		Keywords:       request.Keywords,
		SearchDay:      request.SearchDay,
		SortBy:         request.SortBy,
		PhoneID:        request.PhoneID,
		Search:         request.Search,
	}

	result, err := rtr.engines.SearchListings(searchRequest)
	return listingsSearchResult{Result: result, Error: NewAPIError(err), Message: populateSearchMessage(len(result), request.Keywords)}, err
}

func populateSearchMessage(numberOfResults int, keywords string) string {
	if numberOfResults > 0 {
		return ""
	}

	if numberOfResults == 0 && keywords != "" {
		return "no result found"
	} else if numberOfResults == 0 && keywords == "" {
		return "no deals"
	}

	return ""
}

func (r listingsSearchEndpoint) Validate(request interface{}) error {
	req := request.(listingsSearchRequest)

	for _, listing := range req.ListingTypes {
		isValidListing := false
		for _, definedListings := range shared.ListingTypes {
			if strings.TrimSpace(strings.ToLower((listing))) == definedListings {
				isValidListing = true
				break
			}
		}
		if !isValidListing {
			return helper.ValidationError{Message: fmt.Sprint("listing search failed, invalid 'listingTypes'")}
		}
	}

	if req.SearchDay != "" {
		/*if strings.ToLower(req.SearchDay) != shared.SearchToday &&
			strings.ToLower(req.SearchDay) != shared.SearchTomorrow &&
			strings.ToLower(req.SearchDay) != shared.SearchThisWeek &&
			strings.ToLower(req.SearchDay) != shared.SearchNextWeek {
			return helper.ValidationError{Message: fmt.Sprint("listing search failed, invalid 'searchDay'")}
		}*/
	}

	if strings.TrimSpace(req.SortBy) != "" &&
		strings.ToLower(req.SortBy) != shared.SortByDistance &&
		strings.ToLower(req.SortBy) != shared.SortByTimeLeft &&
		strings.ToLower(req.SortBy) != shared.SortByPrice &&
		req.SortBy != shared.SortByDateAdded &&
		strings.ToLower(req.SortBy) != shared.SortByMostPopular {
		return helper.ValidationError{Message: fmt.Sprint("listing search failed, invalid 'sortBy'")}
	}

	if strings.TrimSpace(req.Location) == "" && (req.Latitude == 0 || req.Longitude == 0) {
		return helper.ValidationError{Message: fmt.Sprint("listings search failed, please provide 'location' or 'latitude' and 'longitude'")}
	}

	for _, dietary := range req.DietaryFilters {
		if strings.ToLower(dietary) != "vegetarian" && strings.ToLower(dietary) != "vegan" && strings.ToLower(dietary) != "gluten free" {
			return helper.ValidationError{Message: fmt.Sprint("listings search failed, invalid dietaryFilters")}
		}
	}

	return nil
}

func (r listingsSearchEndpoint) GetPath() string {
	return "/listings/search"
}

func (r listingsSearchEndpoint) HTTPRequest() interface{} {
	return listingsSearchRequest{}
}
