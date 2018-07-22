package controller

import (
	"context"
	"net/url"
	"strconv"
)

type (
	userEndpoint struct{}
)

var userGet getEndPoint = userEndpoint{}

func (r userEndpoint) Do(ctx context.Context, rtr *router, values url.Values) (interface{}, error) {
	userID, err := strconv.Atoi(values.Get("userId"))
	if err != nil {
		return nil, err
	}

	userInfo, err := rtr.engines.UserGet(userID)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (r userEndpoint) GetPath() string {
	return "/user"
}
