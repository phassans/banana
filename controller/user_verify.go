package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/goware/emailx"
	"github.com/phassans/banana/helper"
)

type (
	verifyUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
		FullName string `json:"fullName"`

		UserID int `json:"userId,omitempty"`
	}

	verifyUserResponse struct {
		verifyUserRequest
		Error *APIError `json:"error,omitempty"`
	}

	verifyUserEndpoint struct{}
)

var userVerify postEndpoint = verifyUserEndpoint{}

func (r verifyUserEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(verifyUserRequest)

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	userInfo, err := rtr.engines.UserVerify(request.Email, request.Password)

	// set values in request object
	request.Password = ""
	request.FullName = userInfo.Name
	request.UserID = userInfo.UserID

	result := verifyUserResponse{verifyUserRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r verifyUserEndpoint) Validate(request interface{}) error {
	input := request.(verifyUserRequest)
	if strings.TrimSpace(input.Email) == "" ||
		strings.TrimSpace(input.Password) == "" {
		return helper.ValidationError{Message: fmt.Sprint("verify user failed, missing fields")}
	}

	if err := emailx.Validate(input.Email); err != nil {
		return helper.ValidationError{Message: fmt.Sprint("verify user failed, invalid email")}
	}

	return nil
}

func (r verifyUserEndpoint) GetPath() string {
	return "/user/verify"
}

func (r verifyUserEndpoint) HTTPRequest() interface{} {
	return verifyUserRequest{}
}
