package controller

import (
	"context"
	"net/url"
)

type allBusinessEndpoint struct{}

var businessAll getEndPoint = allBusinessEndpoint{}

func (r allBusinessEndpoint) Do(ctx context.Context, rtr *router, values url.Values) (interface{}, error) {
	return rtr.engines.GetAllBusiness()
}

func (r allBusinessEndpoint) GetPath() string {
	return "/business/all"
}
