package helper

import (
	"encoding/json"
	"fmt"
)

type ValidationError struct {
	Message string `json:"message,omitempty"`
}

func (v ValidationError) Error() string {
	b, _ := json.Marshal(v)
	return fmt.Sprintf("validation failed with error: %s", string(b))
}

type UserError struct {
	Message string `json:"message,omitempty"`
}

func (v UserError) Error() string {
	b, _ := json.Marshal(v)
	return fmt.Sprintf("user error: %s", string(b))
}

type DuplicateEntity struct {
	BusinessName string `json:"name,omitempty"`
}

func (e DuplicateEntity) Error() string {
	b, _ := json.Marshal(e)
	return fmt.Sprintf("duplicate entity: %s", string(b))
}

type DatabaseError struct {
	DBError string `json:"dbError,omitempty"`
}

func (e DatabaseError) Error() string {
	b, _ := json.Marshal(e)
	return fmt.Sprintf("database error: %s", string(b))
}

type BusinessDoesNotExist struct {
	BusinessName string `json:"name,omitempty"`
}

func (e BusinessDoesNotExist) Error() string {
	b, _ := json.Marshal(e)
	return fmt.Sprintf("business does not exist: %s", string(b))
}
