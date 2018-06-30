package controller

import (
	"context"
	"fmt"

	"github.com/pshassans/banana/helper"
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

var addUser postEndpoint = addUserEndpoint{}

func (r addUserEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(addUserRequest)
	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	err := rtr.engines.AddUser(request.Name, request.Email, request.Password, request.Phone)
	result := addUserResponse{addUserRequest: request, Error: NewAPIError(err)}
	return result, nil
}

func (r addUserEndpoint) Validate(request interface{}) error {
	input := request.(addUserRequest)
	if input.Name == "" || input.Email == "" || input.Password == "" {
		return helper.ValidationError{Message: fmt.Sprint("Add user failed, missing fields")}
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
