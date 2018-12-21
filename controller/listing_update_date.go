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
	listingUpdateDateEndpoint struct{}
)

var listingUpdateDate getEndPoint = listingUpdateDateEndpoint{}

func (r listingUpdateDateEndpoint) Do(ctx context.Context, rtr *router, values url.Values) (interface{}, error) {

	if values.Get("listingId") == "" || values.Get("listingId") == shared.Undefined {
		return nil, helper.ValidationError{Message: fmt.Sprint("user get failed, missing userId")}
	}

	listingID, err := strconv.Atoi(values.Get("listingId"))
	if err != nil {
		return nil, err
	}

	logger := shared.GetLogger()
	logger = logger.With().
		Str("endpoint", r.GetPath()).
		Int("listingId", listingID).Logger()
	logger.Info().Msgf("listing get request")

	err = rtr.engines.UpdateListingDate(listingID)
	if err != nil {
		return nil, err
	}
	return listingInfo, nil
}

func (r listingUpdateDateEndpoint) GetPath() string {
	return "/listing/update/date"
}
