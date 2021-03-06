package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/phassans/banana/helper"
)

// APIError is a HTTP result error.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewAPIError returns a new result error.
func NewAPIError(err error) *APIError {
	if err == nil {
		return nil
	}

	return &APIError{Code: GetErrorStatus(err), Message: err.Error()}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api error [%d]: %s", e.Code, e.Error())
}

// Send writes result error into response.
func (e *APIError) Send(w http.ResponseWriter) error {
	w.WriteHeader(e.Code)
	return json.NewEncoder(w).Encode(e)
}

// GetErrorStatus returns the HTTP status code for a given error type.
func GetErrorStatus(err error) int {
	if e, ok := err.(*APIError); err == nil || ok && e == nil {
		return http.StatusOK
	}
	switch err := err.(type) {
	case helper.LocationError:
		return http.StatusBadRequest
	case helper.DuplicateEntity:
		return http.StatusBadRequest
	case helper.ValidationError:
		return http.StatusBadRequest
	case helper.UserError:
		return http.StatusBadRequest
	case helper.ListingDoesNotExist:
		return http.StatusBadRequest
	case helper.DatabaseError:
		return http.StatusBadRequest
	case helper.BusinessError:
		return http.StatusBadRequest
	case *APIError:
		return err.Code
	default:
		return http.StatusInternalServerError
	}
}

// IsHardError designates whether an error is hard enough to open hystrix circuit breaker
func IsHardError(err error) bool {
	switch err.(type) {
	case hystrix.CircuitError:
		return true
	default:
		return false
	}
}
