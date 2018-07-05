package controller

import (
	"context"
	"net/url"
	"strconv"
)

type (
	listingEndpoint struct{}
)

var listing getEndPoint = listingEndpoint{}

func (r listingEndpoint) Do(ctx context.Context, rtr *router, values url.Values) (interface{}, error) {
	listingID, err := strconv.Atoi(values.Get("listingId"))
	if err != nil {
		return nil, err
	}

	listingInfo, err := rtr.engines.GetListingInfo(listingID)
	if err != nil {
		return nil, err
	}
	return listingInfo, nil
}

func (r listingEndpoint) GetPath() string {
	return "/listing"
}
