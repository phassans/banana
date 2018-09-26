package controller

import (
	"context"
	"net/url"
)

type (
	getStatsEndpoint struct{}
)

var getStats getEndPoint = getStatsEndpoint{}

func (r getStatsEndpoint) Do(ctx context.Context, rtr *router, values url.Values) (interface{}, error) {
	stats, err := rtr.engines.GetStats()
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (r getStatsEndpoint) GetPath() string {
	return "/stats"
}
