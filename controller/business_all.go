package controller

import (
	"context"
	"net/url"

	"github.com/rs/xlog"
)

type allBusinessEndpoint struct{}

var businessAll getEndPoint = allBusinessEndpoint{}

func (r allBusinessEndpoint) Do(ctx context.Context, rtr *router, values url.Values) (interface{}, error) {
	xlog.Infof("GET %s query %+v", r.GetPath(), values)

	return rtr.engines.GetAllBusiness()
}

func (r allBusinessEndpoint) GetPath() string {
	return "/business/all"
}
