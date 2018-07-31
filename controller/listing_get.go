package controller

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/phassans/banana/helper"
	"github.com/rs/xlog"
)

type (
	listingEndpoint struct{}
)

var listingInfo getEndPoint = listingEndpoint{}

func (r listingEndpoint) Do(ctx context.Context, rtr *router, values url.Values) (interface{}, error) {
	xlog.Infof("GET %s query %+v", r.GetPath(), values)

	if values.Get("listingId") == "" {
		return nil, helper.ValidationError{Message: fmt.Sprint("user get failed, missing userId")}
	}

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
