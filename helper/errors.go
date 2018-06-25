package helper

import (
	"encoding/json"
	"fmt"
)

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
