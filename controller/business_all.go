package controller

import (
	"context"
	"net/url"
	"strconv"
)

type allBusinessEndpoint struct{}

var businessAll getEndPoint = allBusinessEndpoint{}

func (r allBusinessEndpoint) Do(ctx context.Context, rtr *router, values url.Values) (interface{}, error) {
	userID, err := strconv.Atoi(values.Get("userId"))
	if err != nil {
		return nil, err
	}

	return rtr.engines.GetAllBusiness(userID)
}

func (r allBusinessEndpoint) GetPath() string {
	return "/business/all"
}
