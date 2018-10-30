package helper

import (
	"encoding/json"
	"fmt"
)

// ValidationError ...
type ValidationError struct {
	Message string `json:"message,omitempty"`
}

func (v ValidationError) Error() string {
	b, _ := json.Marshal(v)
	return fmt.Sprintf("validation error: %s", string(b))
}

// UserError ...
type UserError struct {
	Message string `json:"message,omitempty"`
}

func (v UserError) Error() string {
	b, _ := json.Marshal(v)
	return fmt.Sprintf("user error: %s", string(b))
}

// BusinessError ...
type BusinessError struct {
	Message string `json:"message,omitempty"`
}

func (b BusinessError) Error() string {
	busError, _ := json.Marshal(b)
	return fmt.Sprintf("business error: %s", string(busError))
}

// DuplicateEntity ...
type DuplicateEntity struct {
	Name string `json:"name,omitempty"`
	ID   int    `json:"id,omitempty"`
}

func (e DuplicateEntity) Error() string {
	b, _ := json.Marshal(e)
	return fmt.Sprintf("duplicate entity: %s", string(b))
}

// DatabaseError ...
type DatabaseError struct {
	DBError string `json:"dbError,omitempty"`
}

func (e DatabaseError) Error() string {
	b, _ := json.Marshal(e)
	return fmt.Sprintf("database error: %s", string(b))
}

// BusinessDoesNotExist ...
type BusinessDoesNotExist struct {
	BusinessName string `json:"name,omitempty"`
}

func (e BusinessDoesNotExist) Error() string {
	b, _ := json.Marshal(e)
	return fmt.Sprintf("business does not exist: %s", string(b))
}

// ListingDoesNotExist ...
type ListingDoesNotExist struct {
	ListingID int `json:"listingId,omitempty"`
}

func (e ListingDoesNotExist) Error() string {
	b, _ := json.Marshal(e)
	return fmt.Sprintf("listingId does not exist: %s", string(b))
}

// ValidationError ...
type LocationError struct {
	Message string `json:"message,omitempty"`
}

func (l LocationError) Error() string {
	b, _ := json.Marshal(l)
	return fmt.Sprintf("location error: %s", string(b))
}
