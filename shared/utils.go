package shared

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// NewNullString returns a null able sql string
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// GetTimeIn12HourFormat returns string time in 12 hour format
func GetTimeIn12HourFormat(lTime string) (string, error) {
	if lTime == "" {
		return "", nil
	}

	// determine startTime in format
	st, err := time.Parse(TimeLayout24Hour, strings.TrimSuffix(strings.Split(lTime, "T")[1], "Z"))
	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	return st.Format(TimeLayout12Hour), nil
}
