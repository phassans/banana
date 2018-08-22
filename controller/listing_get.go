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

	if values.Get("listingId") == "" || values.Get("listingId") == shared.Undefined {
		return nil, helper.ValidationError{Message: fmt.Sprint("user get failed, missing userId")}
	}

	listingID, err := strconv.Atoi(values.Get("listingId"))
	if err != nil {
		return nil, err
	}

	var listingDateID int
	if values.Get("listingDateId") != "" && values.Get("listingDateId") != shared.Undefined {
		listingDateID, err = strconv.Atoi(values.Get("listingDateId"))
		if err != nil {
			return nil, err
		}
	}

	var phoneID string
	if values.Get("phoneId") != "" && values.Get("phoneId") != shared.Undefined {
		phoneID = values.Get("phoneId")
	}

	logger := shared.GetLogger()
	logger = logger.With().
		Str("endpoint", r.GetPath()).
		Int("listingId", listingID).
		Int("listingDateId", listingDateID).
		Str("phoneID", phoneID).Logger()
	logger.Info().Msgf("listing get request")

	listingInfo, err := rtr.engines.GetListingInfo(listingID, listingDateID, phoneID)
	if err != nil {
		return nil, err
	}
	return listingInfo, nil
}

func (r listingEndpoint) GetPath() string {
	return "/listing"
}
