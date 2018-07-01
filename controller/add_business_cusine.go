package controller

import (
	"context"
	"fmt"
)

type (
	addBusinessCuisineRequest struct {
		BusinessID int      `json:"business_id"`
		Cuisine    []string `json:"cuisine,omitempty"`
	}

	addBusinessCuisineResponse struct {
		addBusinessCuisineRequest
		Error *APIError `json:"error,omitempty"`
	}

	addBusinessCuisineEndpoint struct{}
)

var addBusinessCuisine postEndpoint = addBusinessCuisineEndpoint{}

func (r addBusinessCuisineEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(addBusinessCuisineRequest)
	if err := r.Validate(requestI); err != nil {
		return nil, err
	}
	fmt.Printf("cuisine: %s", request.Cuisine)

	result := addBusinessCuisineResponse{addBusinessCuisineRequest: request, Error: NewAPIError(nil)}
	return result, nil
}

func (r addBusinessCuisineEndpoint) Validate(request interface{}) error {
	return nil
}

func (r addBusinessCuisineEndpoint) GetPath() string {
	return "/business/cuisine/add"
}

func (r addBusinessCuisineEndpoint) HTTPRequest() interface{} {
	return addBusinessCuisineRequest{}
}
