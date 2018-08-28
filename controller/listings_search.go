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
		Str("sortBy", request.SortBy).Logger()
	logger.Info().Msgf("search request")

	// validate input
	if err := r.Validate(request); err != nil {
		return nil, err
	}

	result, err := rtr.engines.SearchListings(
		request.ListingTypes,
		request.Future,
		request.Latitude,
		request.Longitude,
		request.Location,
		request.PriceFilter,
		request.DietaryFilters,
		request.DistanceFilter,
		request.Keywords,
		request.SortBy,
		request.PhoneID,
	)
	return listingsSearchResult{Result: result, Error: NewAPIError(err), Message: populateSearchMessage(len(result), request.Search)}, err
}

func populateSearchMessage(numberOfResults int, isSearch bool) string {
	if numberOfResults > 0 {
		return ""
	}

	if numberOfResults == 0 && isSearch {
		return "no result found"
	} else if numberOfResults == 0 && !isSearch {
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

	if strings.TrimSpace(req.SortBy) != "" && strings.ToLower(req.SortBy) != "distance" &&
		strings.ToLower(req.SortBy) != "timeleft" && strings.ToLower(req.SortBy) != "price" {
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
