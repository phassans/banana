package controller

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/phassans/banana/helper"
)

type (
	businessEndpoint struct{}
)

var businessInfo getEndPoint = businessEndpoint{}

func (r businessEndpoint) Do(ctx context.Context, rtr *router, values url.Values) (interface{}, error) {

	if values.Get("businessId") == "" {
		return nil, helper.ValidationError{Message: fmt.Sprint("business get failed, missing businessId")}
	}

	businessID, err := strconv.Atoi(values.Get("businessId"))
	if err != nil {
		return nil, err
	}

	businessInfo, err := rtr.engines.GetBusinessInfo(businessID)
	if err != nil {
		return nil, err
	}
	return businessInfo, nil
}

func (r businessEndpoint) GetPath() string {
	return "/business"
}
