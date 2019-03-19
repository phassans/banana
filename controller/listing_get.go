package controller

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
)

type (
	listingEndpoint struct{}
)

var listingInfo getEndPoint = listingEndpoint{}

func (r listingEndpoint) Do(ctx context.Context, rtr *router, values url.Values) (interface{}, error) {
	if err := r.Validate(values); err != nil {
		return nil, err
	}

	listingID, err := strconv.Atoi(values.Get("listingId"))
	if err != nil {
		return nil, err
	}

	phoneID := values.Get("phoneId")
	location := values.Get("location")

	var latitude float64
	if values.Get("latitude") != "" {
		latitude, err = strconv.ParseFloat(values.Get("latitude"), 64)
		if err != nil {
			return nil, err
		}
	}

	var longitude float64
	if values.Get("longitude") != "" {
		longitude, err = strconv.ParseFloat(values.Get("longitude"), 64)
		if err != nil {
			return nil, err
		}
	}

	logger := shared.GetLogger()
	logger = logger.With().
		Str("endpoint", r.GetPath()).
		Int("listingId", listingID).
		Float64("latitude", latitude).
		Float64("longitude", longitude).
		Str("phoneID", phoneID).Logger()
	logger.Info().Msgf("listing get request")

	listingInfo, err := rtr.engines.GetListingInfo(listingID, phoneID, latitude, longitude, location)
	if err != nil {
		return nil, err
	}
	return listingInfo, nil
}

func (r listingEndpoint) GetPath() string {
	return "/listing"
}

func (r listingEndpoint) Validate(values url.Values) error {
	if values.Get("listingId") == "" || values.Get("listingId") == shared.Undefined {
		return helper.ValidationError{Message: fmt.Sprint("listing get failed, missing userId")}
	}

	if values.Get("phoneId") == "" && values.Get("phoneId") == shared.Undefined {
		return helper.ValidationError{Message: fmt.Sprint("listing get failed, missing phoneId")}
	}

	/*if values.Get("latitude") != "" && values.Get("latitude") != shared.Undefined {
		return helper.ValidationError{Message: fmt.Sprint("listing get failed, missing latitude")}
	}

	if values.Get("longitude") != "" && values.Get("longitude") != shared.Undefined {
		return helper.ValidationError{Message: fmt.Sprint("listing get failed, missing longitude")}
	}*/

	return nil
}
