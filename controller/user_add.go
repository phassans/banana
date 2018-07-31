package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/goware/emailx"
	"github.com/phassans/banana/helper"
	"github.com/rs/xlog"
)

type (
	addUserRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Phone    string `json:"phone,omitempty"`
	}

	addUserResponse struct {
		addUserRequest
		Error *APIError `json:"error,omitempty"`
	}

	addUserEndpoint struct{}
)

var userAdd postEndpoint = addUserEndpoint{}

func (r addUserEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(addUserRequest)
	xlog.Infof("POST %s query %+v", r.GetPath(), request)

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	err := rtr.engines.UserAdd(request.Name, request.Email, request.Password, request.Phone)
	result := addUserResponse{addUserRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r addUserEndpoint) Validate(request interface{}) error {
	input := request.(addUserRequest)
	if strings.TrimSpace(input.Name) == "" ||
		strings.TrimSpace(input.Email) == "" ||
		strings.TrimSpace(input.Password) == "" {
		return helper.ValidationError{Message: fmt.Sprint("add user failed, missing fields")}
	}

	if len(input.Password) < 6 {
		return helper.ValidationError{Message: fmt.Sprint("add user failed, password should be atleast 6 characters long")}
	}

	if err := emailx.Validate(input.Email); err != nil {
		return helper.ValidationError{Message: fmt.Sprint("add user failed, invalid email")}
	}

	return nil
}

func (r addUserEndpoint) GetPath() string {
	return "/user/add"
}

func (r addUserEndpoint) HTTPRequest() interface{} {
	return addUserRequest{}
}

func (r addUserEndpoint) HTTPResult() interface{} {
	return addUserResponse{}
}
