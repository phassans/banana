package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/banana/shared"

	"github.com/phassans/banana/helper"
)

type (
	fetchUserDataRequest struct {
		UserID    string `json:"userId"`
		Education bool   `json:"education,omitempty"`
		Company   bool   `json:"company,omitempty"`
	}

	fetchUserDataResponse struct {
		fetchUserDataRequest
		Error      *APIError          `json:"error,omitempty"`
		Companies  []shared.Company   `json:"companies,omitempty"`
		Educations []shared.Education `json:"educations,omitempty"`
	}

	fetchUserEducationEndpoint struct{}
)

var fetchUser postEndpoint = fetchUserEducationEndpoint{}

func (r fetchUserEducationEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(fetchUserDataRequest)
	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	var err error
	if request.Company {
		comps, _ := rtr.engines.FetchUserCompanies(request.UserID)
		return fetchUserDataResponse{fetchUserDataRequest: request, Error: NewAPIError(err), Companies: comps}, nil
	} else if request.Education {
		edus, err := rtr.engines.FetchUserEducation(request.UserID)
		return fetchUserDataResponse{fetchUserDataRequest: request, Error: NewAPIError(err), Educations: edus}, nil

	}

	result := fetchUserDataResponse{fetchUserDataRequest: request, Error: NewAPIError(err)}
	return result, nil
}

func (r fetchUserEducationEndpoint) Validate(request interface{}) error {
	input := request.(fetchUserDataRequest)
	if strings.TrimSpace(input.UserID) == "" {
		return helper.ValidationError{Message: fmt.Sprint("fetch user education failed, missing fields")}
	}

	if input.Education == false && input.Company == false {
		return helper.ValidationError{Message: fmt.Sprint("fetch user education failed, missing fields. Either company or education should set to true")}
	}

	return nil
}

func (r fetchUserEducationEndpoint) GetPath() string {
	return "/user/fetch"
}

func (r fetchUserEducationEndpoint) HTTPRequest() interface{} {
	return fetchUserDataRequest{}
}

func (r fetchUserEducationEndpoint) HTTPResult() interface{} {
	return fetchUserDataResponse{}
}
