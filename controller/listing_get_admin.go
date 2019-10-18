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
	listingAdminEndpoint struct{}
)

var listingAdminInfo getEndPoint = listingAdminEndpoint{}

func (r listingAdminEndpoint) Do(ctx context.Context, rtr *router, values url.Values) (interface{}, error) {
	if err := r.Validate(values); err != nil {
		return nil, err
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

	listingInfo, err := rtr.engines.GetListingByIDAdmin(listingID)
	if err != nil {
		return nil, err
	}
	return listingInfo, nil
}

func (r listingAdminEndpoint) GetPath() string {
	return "/listing/admin"
}

func (r listingAdminEndpoint) Validate(values url.Values) error {
	if values.Get("listingId") == "" || values.Get("listingId") == shared.Undefined {
		return helper.ValidationError{Message: fmt.Sprint("listing get failed, missing userId")}
	}

	return nil
}
