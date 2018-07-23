package controller

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/phassans/banana/helper"
)

type (
	userEndpoint struct{}
)

var userGet getEndPoint = userEndpoint{}

func (r userEndpoint) Do(ctx context.Context, rtr *router, values url.Values) (interface{}, error) {

	if values.Get("userId") == "" {
		return nil, helper.ValidationError{Message: fmt.Sprint("user get failed, missing userId")}
	}

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
