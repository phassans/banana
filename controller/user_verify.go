package controller

import (
	"context"
	"fmt"

	"github.com/phassans/banana/helper"
)

type (
	verifyUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	verifyUserResponse struct {
		verifyUserRequest
		UserID int       `json:"user_id"`
		Error  *APIError `json:"error,omitempty"`
	}

	verifyUserEndpoint struct{}
)

var userVerify postEndpoint = verifyUserEndpoint{}

func (r verifyUserEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(verifyUserRequest)
	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	userID, err := rtr.engines.VerifyUser(request.Email, request.Password)
	result := verifyUserResponse{verifyUserRequest: request, UserID: userID, Error: NewAPIError(err)}
	return result, err
}

func (r verifyUserEndpoint) Validate(request interface{}) error {
	input := request.(verifyUserRequest)
	if input.Email == "" || input.Password == "" {
		return helper.ValidationError{Message: fmt.Sprint("Add user failed, missing fields")}
	}

	return nil
}

func (r verifyUserEndpoint) GetPath() string {
	return "/user/verify"
}

func (r verifyUserEndpoint) HTTPRequest() interface{} {
	return verifyUserRequest{}
}
