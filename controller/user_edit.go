package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/goware/emailx"
	"github.com/phassans/banana/helper"
)

type (
	editUserRequest struct {
		UserId   int    `json:"userId"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Phone    string `json:"phone,omitempty"`
	}

	editUserResponse struct {
		editUserRequest
		Error *APIError `json:"error,omitempty"`
	}

	editUserEndpoint struct{}
)

var userEdit postEndpoint = editUserEndpoint{}

func (r editUserEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(editUserRequest)
	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	err := rtr.engines.UserEdit(request.UserId, request.Name, request.Email, request.Password, request.Phone)
	result := editUserResponse{editUserRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r editUserEndpoint) Validate(request interface{}) error {
	input := request.(editUserRequest)
	if strings.TrimSpace(input.Name) == "" ||
		strings.TrimSpace(input.Email) == "" ||
		strings.TrimSpace(input.Password) == "" {
		return helper.ValidationError{Message: fmt.Sprint("edit user failed, missing fields")}
	}

	if len(input.Password) < 6 {
		return helper.ValidationError{Message: fmt.Sprint("add user failed, password should be atleast 6 characters long")}
	}

	if err := emailx.Validate(input.Email); err != nil {
		return helper.ValidationError{Message: fmt.Sprint("edit user failed, invalid email")}
	}

	return nil
}

func (r editUserEndpoint) GetPath() string {
	return "/user/edit"
}

func (r editUserEndpoint) HTTPRequest() interface{} {
	return editUserRequest{}
}

func (r editUserEndpoint) HTTPResult() interface{} {
	return editUserResponse{}
}
