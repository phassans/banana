package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/shared"
)

type (
	reportInaccurateRequest struct {
		PhoneID   string `json:"phoneId"`
		ListingID int    `json:"listingId"`
	}

	reportInaccurateResponse struct {
		reportInaccurateRequest
		Error *APIError `json:"error,omitempty"`
	}

	reportInaccurateEndpoint struct{}
)

var reportInaccurate postEndpoint = reportInaccurateEndpoint{}

func (r reportInaccurateEndpoint) Execute(ctx context.Context, rtr *router, requestI interface{}) (interface{}, error) {
	request := requestI.(reportInaccurateRequest)

	logger := shared.GetLogger()
	logger = logger.With().
		Str("endpoint", r.GetPath()).
		Str("phoneId", request.PhoneID).
		Int("listingId", request.ListingID).Logger()
	logger.Info().Msgf("reportInaccurate add request")

	if err := r.Validate(requestI); err != nil {
		return nil, err
	}

	err := rtr.engines.ReportInaccurate(request.PhoneID, request.ListingID)
	result := reportInaccurateResponse{reportInaccurateRequest: request, Error: NewAPIError(err)}
	return result, err
}

func (r reportInaccurateEndpoint) Validate(request interface{}) error {
	req := request.(reportInaccurateRequest)

	if strings.TrimSpace(req.PhoneID) == "" {
		return helper.ValidationError{Message: fmt.Sprint("favorite add failed, please provide 'phoneId'")}
	}

	if req.ListingID == 0 {
		return helper.ValidationError{Message: fmt.Sprint("favorite add failed, please provide 'listingId'")}
	}

	return nil
}

func (r reportInaccurateEndpoint) GetPath() string {
	return "/reportinaccurate"
}

func (r reportInaccurateEndpoint) HTTPRequest() interface{} {
	return reportInaccurateRequest{}
}
